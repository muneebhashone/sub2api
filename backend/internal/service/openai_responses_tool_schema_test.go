package service

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// issue #5364 的最小复现体：Codex Desktop 内置 automation_update 带
// parameters.type = null，upstream 回 400 invalid_function_parameters。
func TestSanitizeOpenAIResponsesToolParameterTypes_TopLevelFunctionTool(t *testing.T) {
	body := []byte(`{
		"model": "gpt-5.6-sol",
		"input": "Reply with OK.",
		"stream": false,
		"tools": [
			{
				"type": "function",
				"name": "automation_update",
				"description": "Update an automation.",
				"parameters": {"type": null, "properties": {placeholderplaceholder
		placeholder
		]
placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.True(t, changed)
	require.Equal(t, "object", gjson.GetBytes(sanitized, "tools.0.parameters.type").String())
	// 只改 type，工具其余定义原样保留。
	require.Equal(t, "automation_update", gjson.GetBytes(sanitized, "tools.0.name").String())
	require.Equal(t, "Update an automation.", gjson.GetBytes(sanitized, "tools.0.description").String())
	require.True(t, gjson.GetBytes(sanitized, "tools.0.parameters.properties").IsObject())
	// 请求体其余字段不受影响。
	require.Equal(t, "gpt-5.6-sol", gjson.GetBytes(sanitized, "model").String())
	require.Equal(t, "Reply with OK.", gjson.GetBytes(sanitized, "input").String())
placeholder

// 合法 Schema 必须原样返回：changed=false 且字节不变，避免无谓重写打散
// prompt cache 前缀。
func TestSanitizeOpenAIResponsesToolParameterTypes_ValidSchemaUntouched(t *testing.T) {
	body := []byte(`{"tools":[{"type":"function","name":"ok","parameters":{"type":"object","properties":{placeholderplaceholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.False(t, changed)
	require.Equal(t, string(body), string(sanitized))
placeholder

// 缺失 type 的 Schema 本身合法（等价于不约束），不得补写——补写会收窄客户端语义。
func TestSanitizeOpenAIResponsesToolParameterTypes_MissingTypeNotInvented(t *testing.T) {
	body := []byte(`{"tools":[{"type":"function","name":"ok","parameters":{"properties":{placeholderplaceholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.False(t, changed)
	require.False(t, gjson.GetBytes(sanitized, "tools.0.parameters.type").Exists())
placeholder

func TestSanitizeOpenAIResponsesToolParameterTypes_ObjectOnlyRootUnion(t *testing.T) {
	body := []byte(`{"tools":[
		{"type":"function","name":"automation_update","parameters":{"oneOf":[{"type":"object","properties":{"id":{"type":"string"placeholderplaceholderplaceholder,{"anyOf":[{"type":"object"placeholder,{"type":"object","properties":{placeholderplaceholder]placeholder]placeholderplaceholder,
		{"type":"function","name":"codex_app__automation_update","parameters":{"anyOf":[{"type":"object"placeholder]placeholderplaceholder,
		{"type":"function","name":"unrelated_tool","parameters":{"oneOf":[{"type":"object"placeholder]placeholderplaceholder
	]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.True(t, changed)
	for i := 0; i < 3; i++ {
		require.Equal(t, "object", gjson.GetBytes(sanitized, "tools."+strconv.Itoa(i)+".parameters.type").String())
placeholder
	require.Equal(t, "string", gjson.GetBytes(sanitized, "tools.0.parameters.oneOf.0.properties.id.type").String())
placeholder

func TestSanitizeOpenAIResponsesToolParameterTypes_EscapedUnionKeyword(t *testing.T) {
	body := []byte(`{"tools":[{"type":"function","name":"other","parameters":{"one\u004ff":[{"type":"object"placeholder]placeholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.True(t, changed)
	require.Equal(t, "object", gjson.GetBytes(sanitized, "tools.0.parameters.type").String())
placeholder

func TestSanitizeOpenAIResponsesToolParameterTypes_ObjectUnionMissingTypeSafetyBoundary(t *testing.T) {
	tests := []struct {
		name string
		body string
placeholder{
		{name: "mixed union", body: `{"tools":[{"type":"function","name":"automation_update","parameters":{"oneOf":[{"type":"object"placeholder,{"type":"string"placeholder]placeholderplaceholder]placeholder`placeholder,
		{name: "unproven branch", body: `{"tools":[{"type":"function","name":"automation_update","parameters":{"oneOf":[{"properties":{placeholderplaceholder]placeholderplaceholder]placeholder`placeholder,
		{name: "empty union", body: `{"tools":[{"type":"function","name":"automation_update","parameters":{"oneOf":[]placeholderplaceholder]placeholder`placeholder,
		{name: "existing null handled by existing repair", body: `{"tools":[{"type":"function","name":"automation_update","parameters":{"type":null,"oneOf":[{"type":"object"placeholder]placeholderplaceholder]placeholder`placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes([]byte(tt.body))
		placeholder
			if tt.name == "existing null handled by existing repair" {
				require.True(t, changed)
				require.Equal(t, "object", gjson.GetBytes(sanitized, "tools.0.parameters.type").String())
				return
		placeholder
			require.False(t, changed)
			require.Equal(t, tt.body, string(sanitized))
	placeholder)
placeholder
placeholder

// 多轮历史：工具定义沉进 input 后，upstream 报错路径形如
// input[N].tools[i].tools[j].parameters，两层都要修。
func TestSanitizeOpenAIResponsesToolParameterTypes_NestedHistoryTools(t *testing.T) {
	body := []byte(`{
		"input": [
			{"type": "message", "role": "user", "content": "hi"placeholder,
			{
				"type": "additional_tools",
				"role": "developer",
				"parameters": {"type": nullplaceholder,
				"tools": [
					{
						"type": "namespace",
						"name": "codex_app",
						"tools": [
							{"type": "function", "name": "noop", "parameters": {"type": "object"placeholderplaceholder,
							{"type": "function", "name": "automation_update", "parameters": {"type": nullplaceholderplaceholder
						]
				placeholder,
					{"type": "function", "name": "outer", "parameters": {"type": nullplaceholderplaceholder
				]
		placeholder
		]
placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.True(t, changed)
	require.Equal(t, "object", gjson.GetBytes(sanitized, "input.1.tools.0.tools.1.parameters.type").String())
	require.Equal(t, "object", gjson.GetBytes(sanitized, "input.1.tools.1.parameters.type").String())
	require.Equal(t, "null", gjson.GetBytes(sanitized, "input.1.parameters.type").Raw)
	// 原本合法的兄弟条目保持不变。
	require.Equal(t, "object", gjson.GetBytes(sanitized, "input.1.tools.0.tools.0.parameters.type").String())
	require.Equal(t, "hi", gjson.GetBytes(sanitized, "input.0.content").String())
placeholder

// ChatCompletions 形态的工具（{type:"function", function:{...placeholderplaceholder）同样可能出现在
// Responses 请求里，见 normalizeCodexTools。
func TestSanitizeOpenAIResponsesToolParameterTypes_ChatCompletionsShape(t *testing.T) {
	body := []byte(`{"tools":[{"type":"function","function":{"name":"legacy","parameters":{"type":nullplaceholderplaceholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.True(t, changed)
	require.Equal(t, "object", gjson.GetBytes(sanitized, "tools.0.function.parameters.type").String())
	require.Equal(t, "legacy", gjson.GetBytes(sanitized, "tools.0.function.name").String())
placeholder

func TestSanitizeOpenAIResponsesToolParameterTypes_InputItemShapes(t *testing.T) {
	body := []byte(`{"input":[{"parameters":{"type":nullplaceholder,"name":"direct","type":"function"placeholder,{"function":{"name":"wrapped","parameters":{"type":nullplaceholderplaceholder,"type":"function"placeholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.True(t, changed)
	require.Equal(t, "object", gjson.GetBytes(sanitized, "input.0.parameters.type").String())
	require.Equal(t, "object", gjson.GetBytes(sanitized, "input.1.function.parameters.type").String())
placeholder

func TestSanitizeOpenAIResponsesToolParameterTypes_DoesNotTreatMessagesAsTools(t *testing.T) {
	body := []byte(`{"input":[{"type":"message","role":"user","parameters":{"type":nullplaceholder,"function":{"parameters":{"type":nullplaceholderplaceholder,"metadata":{"parameters":{"type":nullplaceholderplaceholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.False(t, changed)
	require.Equal(t, string(body), string(sanitized))
placeholder

func TestSanitizeOpenAIResponsesToolParameterTypes_PreservesUntouchedBytes(t *testing.T) {
	body := []byte("{\n  \"z\":1e+999, \"integer\":900719925474099312345, \"escaped\":\"\\u003chtml\\u003e\", \"tools\":[{\"parameters\":{\"properties\":{placeholder, \"type\" : nullplaceholder,\"name\":\"x\"placeholder], \"a\":-0.00\nplaceholder")
	want := []byte("{\n  \"z\":1e+999, \"integer\":900719925474099312345, \"escaped\":\"\\u003chtml\\u003e\", \"tools\":[{\"parameters\":{\"properties\":{placeholder, \"type\" : \"object\"placeholder,\"name\":\"x\"placeholder], \"a\":-0.00\nplaceholder")

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.True(t, changed)
	require.Equal(t, string(want), string(sanitized))
placeholder

func TestSanitizeOpenAIResponsesToolParameterTypes_EscapedKey(t *testing.T) {
	body := []byte(`{"tools":[{"type":"function","parameters":{"ty\u0070e":nullplaceholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.True(t, changed)
	require.Equal(t, `{"tools":[{"type":"function","parameters":{"ty\u0070e":"object"placeholderplaceholder]placeholder`, string(sanitized))
placeholder

func TestSanitizeOpenAIResponsesToolParameterTypes_DuplicateTypeIsAmbiguousNoOp(t *testing.T) {
	for _, body := range [][]byte{
		[]byte(`{"tools":[{"type":"function","parameters":{"type":null,"type":"string"placeholderplaceholder]placeholder`),
		[]byte(`{"tools":[{"type":"function","parameters":{"type":"string","ty\u0070e":nullplaceholderplaceholder]placeholder`),
		[]byte(`{"tools":[{"type":"function","parameters":{"type":null,"type":nullplaceholderplaceholder]placeholder`),
placeholder {
		sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)
	placeholder
		require.False(t, changed)
		require.Equal(t, string(body), string(sanitized))
placeholder
placeholder

func TestSanitizeOpenAIResponsesToolSchemas_InvalidAndTrailingJSON(t *testing.T) {
	for _, tc := range []struct {
		name string
		body []byte
		call func([]byte) ([]byte, bool, error)
placeholder{
		{"null_type_trailing", []byte(`{"tools":[{"parameters":{"type":nullplaceholderplaceholder]placeholder trailing`), sanitizeOpenAIResponsesToolParameterTypesplaceholder,
		{"pattern_trailing", []byte(`{"tools":[{"parameters":{"pattern":"(?=x)"placeholderplaceholder]placeholder trailing`), sanitizeOpenAIResponsesToolSchemaPatternsplaceholder,
		{"invalid_utf8", append([]byte(`{"tools":[{"parameters":{"pattern":"(?=x)"placeholderplaceholder],"bad":"`), 0xff, '"', 'placeholder'), sanitizeOpenAIResponsesToolSchemaPatternsplaceholder,
placeholder {
		t.Run(tc.name, func(t *testing.T) {
			_, changed, err := tc.call(tc.body)
		placeholder
			require.False(t, changed)
	placeholder)
placeholder
placeholder

func TestOpenAIResponsesToolSchemaCapabilities_PlatformBoundary(t *testing.T) {
	tests := []struct {
		platform         string
		repairNullType   bool
		removeLookaround bool
placeholder{
		{PlatformOpenAI, true, trueplaceholder,
		{PlatformAnthropic, true, falseplaceholder,
		{PlatformKimi, true, falseplaceholder,
		{PlatformZhipu, true, falseplaceholder,
		{PlatformDeepseek, true, falseplaceholder,
		{PlatformGrok, false, falseplaceholder,
		{PlatformGemini, false, falseplaceholder,
		{PlatformAntigravity, false, falseplaceholder,
		{PlatformComposite, false, falseplaceholder,
		{"", false, falseplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			require.Equal(t, tt.repairNullType, shouldRepairOpenAIResponsesNullToolSchemaType(tt.platform))
			require.Equal(t, tt.removeLookaround, shouldSanitizeOpenAIResponsesToolSchemaPatterns(tt.platform))
	placeholder)
placeholder
placeholder

func TestSanitizeOpenAIResponsesToolSchemasForPlatform_ReplayBoundary(t *testing.T) {
	body := []byte(`{"tools":[{"type":"function","parameters":{"type":null,"properties":{"query":{"type":"string","pattern":"(?=keep)"placeholderplaceholderplaceholderplaceholder]placeholder`)

	// A malformed tool definition may be replayed after account failover. Every
	// compatible account must repair it, while non-OpenAI providers retain their
	// supported regex semantics.
	for _, platform := range []string{PlatformAnthropic, PlatformKimi, PlatformZhipu, PlatformDeepseekplaceholder {
		t.Run(platform, func(t *testing.T) {
			for attempt := 0; attempt < 2; attempt++ {
				normalized, changed, err := sanitizeOpenAIResponsesToolSchemasForPlatform(body, platform)
			placeholder
				require.True(t, changed)
				require.Equal(t, "object", gjson.GetBytes(normalized, "tools.0.parameters.type").String())
				require.Equal(t, "(?=keep)", gjson.GetBytes(normalized, "tools.0.parameters.properties.query.pattern").String())
		placeholder
	placeholder)
placeholder

	openAI, changed, err := sanitizeOpenAIResponsesToolSchemasForPlatform(body, PlatformOpenAI)
placeholder
	require.True(t, changed)
	require.Equal(t, "object", gjson.GetBytes(openAI, "tools.0.parameters.type").String())
	require.False(t, gjson.GetBytes(openAI, "tools.0.parameters.properties.query.pattern").Exists())

	unsupported, changed, err := sanitizeOpenAIResponsesToolSchemasForPlatform(body, PlatformGrok)
placeholder
	require.False(t, changed)
	require.Equal(t, string(body), string(unsupported))
placeholder

// 索引映射：只有坏条目被改，前后兄弟条目按原下标保持不变。
func TestSanitizeOpenAIResponsesToolParameterTypes_OnlyOffendingIndexRewritten(t *testing.T) {
	body := []byte(`{"tools":[
		{"type":"function","name":"a","parameters":{"type":"object"placeholderplaceholder,
		{"type":"function","name":"b","parameters":{"type":nullplaceholderplaceholder,
		{"type":"function","name":"c","parameters":{"type":"object"placeholderplaceholder
	]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.True(t, changed)
	require.Equal(t, "a", gjson.GetBytes(sanitized, "tools.0.name").String())
	require.Equal(t, "b", gjson.GetBytes(sanitized, "tools.1.name").String())
	require.Equal(t, "c", gjson.GetBytes(sanitized, "tools.2.name").String())
	require.Equal(t, "object", gjson.GetBytes(sanitized, "tools.1.parameters.type").String())
	require.Equal(t, 3, int(gjson.GetBytes(sanitized, "tools.#").Int()))
placeholder

// 畸形/非常规形态不得 panic，且一律按不变处理。
func TestSanitizeOpenAIResponsesToolParameterTypes_MalformedShapesAreNoOps(t *testing.T) {
	cases := []struct {
		name string
		body string
placeholder{
		{"empty_body", ``placeholder,
		{"no_tools", `{"model":"gpt-5.6-sol","input":"hi"placeholder`placeholder,
		{"tools_null", `{"tools":nullplaceholder`placeholder,
		{"tools_object", `{"tools":{"type":"function"placeholderplaceholder`placeholder,
		{"tool_is_string", `{"tools":["freeform"]placeholder`placeholder,
		{"parameters_is_string", `{"tools":[{"type":"function","parameters":"nope"placeholder]placeholder`placeholder,
		{"parameters_null", `{"tools":[{"type":"function","parameters":nullplaceholder]placeholder`placeholder,
		{"input_string", `{"input":"hi","tools":[]placeholder`placeholder,
		{"input_item_not_object", `{"input":["hi"]placeholder`placeholder,
		{"type_already_array", `{"tools":[{"type":"function","parameters":{"type":["object","null"]placeholderplaceholder]placeholder`placeholder,
placeholder

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes([]byte(tc.body))
		placeholder
			require.False(t, changed)
			require.Equal(t, tc.body, string(sanitized))
	placeholder)
placeholder
placeholder

// 递归深度守卫：超深嵌套只做截断，不递归到栈溢出，也不报错。
func TestSanitizeOpenAIResponsesToolParameterTypes_DepthGuard(t *testing.T) {
	tool := map[string]any{"type": "function", "name": "deep", "parameters": map[string]any{"type": nilplaceholderplaceholder
	for i := 0; i < 12; i++ {
		tool = map[string]any{"type": "namespace", "tools": []any{toolplaceholderplaceholder
placeholder
	body, err := json.Marshal(map[string]any{"tools": []any{toolplaceholderplaceholder)
placeholder

	require.NotPanics(t, func() {
		_, _, sanitizeErr := sanitizeOpenAIResponsesToolParameterTypes(body)
		require.NoError(t, sanitizeErr)
placeholder)
placeholder

// 输出必须是合法 JSON，且除目标字段外与输入等价。
func TestSanitizeOpenAIResponsesToolParameterTypes_OutputStaysValidJSON(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","tool_choice":"none","store":false,"tools":[{"type":"function","name":"automation_update","parameters":{"type":null,"properties":{placeholderplaceholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)
placeholder
	require.True(t, changed)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(sanitized, &decoded))
	require.Equal(t, "gpt-5.5", decoded["model"])
	require.Equal(t, "none", decoded["tool_choice"])
	require.Equal(t, false, decoded["store"])
placeholder

// 输入 body 是调用方持有的缓冲区（Forward 里 canonicalImageIntentBody 与它同源），
// 净化必须返回新切片，绝不能就地改写。
func TestSanitizeOpenAIResponsesToolParameterTypes_DoesNotMutateInputBody(t *testing.T) {
	body := []byte(`{"model":"gpt-5.6-sol","tools":[{"type":"function","name":"a","parameters":{"type":nullplaceholderplaceholder]placeholder`)
	original := append([]byte(nil), body...)

	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(body)

placeholder
	require.True(t, changed)
	require.Equal(t, string(original), string(body), "调用方的 body 不得被就地改写")
	require.NotEqual(t, string(original), string(sanitized))
placeholder

func TestSanitizeOpenAIResponsesToolSchemaPatterns_RemovesLookaroundOnly(t *testing.T) {
	body := []byte(`{"tools":[{"type":"function","name":"search","parameters":{"type":"object","properties":{"q":{"type":"string","pattern":"^(?=.*foo)[a-z]+$"placeholder,"id":{"type":"string","pattern":"^[a-z]+$"placeholder,"z":{"type":"string","pattern":"(?<!bad)ok"placeholderplaceholderplaceholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolSchemaPatterns(body)
placeholder
	require.True(t, changed)
	require.False(t, gjson.GetBytes(sanitized, "tools.0.parameters.properties.q.pattern").Exists())
	require.False(t, gjson.GetBytes(sanitized, "tools.0.parameters.properties.z.pattern").Exists())
	require.Equal(t, "^[a-z]+$", gjson.GetBytes(sanitized, "tools.0.parameters.properties.id.pattern").String())
placeholder

func TestSanitizeOpenAIResponsesToolSchemaPatterns_PreservesRawDocumentRepresentation(t *testing.T) {
	body := []byte("{\n \"number\":1e+999,\n \"tools\":[{\"parameters\":{\"pattern\":\"(?=drop)\", \"description\":\"\\u003ctag\\u003e\", \"type\":\"string\"placeholderplaceholder],\n \"tail\":-0.00\nplaceholder")
	want := []byte("{\n \"number\":1e+999,\n \"tools\":[{\"parameters\":{ \"description\":\"\\u003ctag\\u003e\", \"type\":\"string\"placeholderplaceholder],\n \"tail\":-0.00\nplaceholder")

	sanitized, changed, err := sanitizeOpenAIResponsesToolSchemaPatterns(body)

placeholder
	require.True(t, changed)
	require.Equal(t, string(want), string(sanitized))
placeholder

func TestSanitizeOpenAIResponsesToolSchemaPatterns_AllMemberPositions(t *testing.T) {
	body := []byte(`{"tools":[{"parameters":{"properties":{"first":{"pattern":"(?=a)","type":"string"placeholder,"middle":{"title":"x","pattern":"(?!b)","type":"string"placeholder,"last":{"type":"string","pattern":"(?<=c)"placeholder,"only":{"pattern":"(?<!d)"placeholderplaceholderplaceholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolSchemaPatterns(body)

placeholder
	require.True(t, changed)
	for _, name := range []string{"first", "middle", "last", "only"placeholder {
		require.False(t, gjson.GetBytes(sanitized, "tools.0.parameters.properties."+name+".pattern").Exists(), name)
		require.True(t, gjson.GetBytes(sanitized, "tools.0.parameters.properties."+name).IsObject(), name)
placeholder
placeholder

func TestSanitizeOpenAIResponsesToolSchemaPatterns_TraversesConservativeSchemaKeywords(t *testing.T) {
	body := []byte(`{"input":[{"type":"function","name":"direct","parameters":{"type":"object","properties":{"p":{"anyOf":[{"pattern":"(?=any)"placeholder],"items":{"pattern":"(?!item)"placeholderplaceholderplaceholder,"$defs":{"d":{"pattern":"(?<=def)"placeholderplaceholder,"additionalProperties":{"pattern":"(?<!additional)"placeholderplaceholderplaceholder,{"type":"function","function":{"name":"wrapped","parameters":{"oneOf":[{"pattern":"(?=wrapped)"placeholder]placeholderplaceholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolSchemaPatterns(body)

placeholder
	require.True(t, changed)
	for _, path := range []string{
		"input.0.parameters.properties.p.anyOf.0.pattern",
		"input.0.parameters.properties.p.items.pattern",
		"input.0.parameters.$defs.d.pattern",
		"input.0.parameters.additionalProperties.pattern",
		"input.1.function.parameters.oneOf.0.pattern",
placeholder {
		require.False(t, gjson.GetBytes(sanitized, path).Exists(), path)
placeholder
placeholder

func TestSanitizeOpenAIResponsesToolSchemaPatterns_DoesNotTraverseMessageMetadata(t *testing.T) {
	body := []byte(`{"input":[{"parameters":{"pattern":"(?=missing-type)"placeholderplaceholder,{"type":"message","parameters":{"pattern":"(?=parameters)"placeholder,"function":{"parameters":{"pattern":"(?=function)"placeholderplaceholder,"metadata":{"tools":[{"type":"function","parameters":{"pattern":"(?=metadata)"placeholderplaceholder]placeholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolSchemaPatterns(body)

placeholder
	require.False(t, changed)
	require.Equal(t, string(body), string(sanitized))
placeholder

func TestSanitizeOpenAIResponsesToolSchemaPatterns_DecodesEscapedLookaround(t *testing.T) {
	for _, pattern := range []string{`\u0028?\u003dfoo`, `(\u003f\u003dfoo)`placeholder {
		body := []byte(`{"tools":[{"parameters":{"type":"string","pattern":"` + pattern + `"placeholderplaceholder]placeholder`)
		sanitized, changed, err := sanitizeOpenAIResponsesToolSchemaPatterns(body)
	placeholder
		require.True(t, changed)
		require.False(t, gjson.GetBytes(sanitized, "tools.0.parameters.pattern").Exists())
placeholder
placeholder

func TestSanitizeOpenAIResponsesToolSchemaPatterns_EscapedKey(t *testing.T) {
	body := []byte(`{"to\u006fls":[{"par\u0061meters":{"type":"string","pat\u0074ern":"(?=drop)"placeholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolSchemaPatterns(body)

placeholder
	require.True(t, changed)
	require.Equal(t, `{"to\u006fls":[{"par\u0061meters":{"type":"string"placeholderplaceholder]placeholder`, string(sanitized))
placeholder

func TestSanitizeOpenAIResponsesToolSchemaPatterns_DuplicateKeysStayValid(t *testing.T) {
	for _, tc := range []struct {
		body string
		want string
placeholder{
		{
			`{"tools":[{"parameters":{"pattern":"(?=drop)","pattern":"^[a-z]+$","type":"string"placeholderplaceholder]placeholder`,
			`{"tools":[{"parameters":{"pattern":"^[a-z]+$","type":"string"placeholderplaceholder]placeholder`,
	placeholder,
		{
			`{"tools":[{"parameters":{"type":"string","pattern":"(?=first)","pattern":"(?!second)"placeholderplaceholder]placeholder`,
			`{"tools":[{"parameters":{"type":"string"placeholderplaceholder]placeholder`,
	placeholder,
placeholder {
		sanitized, changed, err := sanitizeOpenAIResponsesToolSchemaPatterns([]byte(tc.body))
	placeholder
		require.True(t, changed)
		require.True(t, json.Valid(sanitized))
		require.Equal(t, tc.want, string(sanitized))
placeholder
placeholder

func TestOpenAIResponsesToolSchemaPlatformGate_APIKeyAndOAuth(t *testing.T) {
	body := []byte(`{"tools":[{"type":"function","parameters":{"type":null,"pattern":"(?=drop)"placeholderplaceholder]placeholder`)
	for _, accountType := range []string{AccountTypeAPIKey, AccountTypeOAuthplaceholder {
		t.Run(accountType, func(t *testing.T) {
			normalized, changed, err := normalizeOpenAIResponsesWebSocketCompatibilityBody(body, &Account{
				Platform: PlatformOpenAI,
				Type:     accountType,
		placeholder)
		placeholder
			require.True(t, changed)
			require.Equal(t, "object", gjson.GetBytes(normalized, "tools.0.parameters.type").String())
			require.False(t, gjson.GetBytes(normalized, "tools.0.parameters.pattern").Exists())
	placeholder)
placeholder

	normalized, changed, err := normalizeOpenAIResponsesWebSocketCompatibilityBody(body, &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeAPIKey,
placeholder)
placeholder
	require.False(t, changed)
	require.Equal(t, string(body), string(normalized))
placeholder

func TestSanitizeOpenAIResponsesToolSchemas_DepthLimitReturnsOriginal(t *testing.T) {
	body := []byte(`{"tools":[{"parameters":`)
	for i := 0; i < openAIResponsesToolSchemaMaxDepth+2; i++ {
		body = append(body, []byte(`{"allOf":[`)...)
placeholder
	body = append(body, []byte(`{"pattern":"(?=too-deep)"placeholder`)...)
	for i := 0; i < openAIResponsesToolSchemaMaxDepth+2; i++ {
		body = append(body, []byte(`]placeholder`)...)
placeholder
	body = append(body, []byte(`placeholder]placeholder`)...)
	require.True(t, json.Valid(body))

	sanitized, changed, err := sanitizeOpenAIResponsesToolSchemaPatterns(body)

placeholder
	require.False(t, changed)
	require.Equal(t, string(body), string(sanitized))
placeholder

func TestSanitizeOpenAIResponsesToolSchemaPatterns_DoesNotTouchUserInputPattern(t *testing.T) {
	body := []byte(`{"input":{"pattern":"(?=keep)"placeholder,"metadata":{"pattern":"(?!keep)"placeholderplaceholder`)
	sanitized, changed, err := sanitizeOpenAIResponsesToolSchemaPatterns(body)
placeholder
	require.False(t, changed)
	require.Equal(t, string(body), string(sanitized))
placeholder

func TestSanitizeOpenAIResponsesToolSchemaPatterns_DoesNotTraverseInstanceData(t *testing.T) {
	body := []byte(`{"tools":[{"type":"function","parameters":{"type":"object","properties":{"config":{"type":"object","pattern":"(?=remove)","default":{"pattern":"(?=keep-default)"placeholder,"examples":[{"pattern":"(?=keep-example)"placeholder],"const":{"pattern":"(?=keep-const)"placeholder,"enum":[{"pattern":"(?=keep-enum)"placeholder]placeholderplaceholderplaceholderplaceholder]placeholder`)

	sanitized, changed, err := sanitizeOpenAIResponsesToolSchemaPatterns(body)
placeholder
	require.True(t, changed)
	require.False(t, gjson.GetBytes(sanitized, "tools.0.parameters.properties.config.pattern").Exists())
	require.Equal(t, "(?=keep-default)", gjson.GetBytes(sanitized, "tools.0.parameters.properties.config.default.pattern").String())
	require.Equal(t, "(?=keep-example)", gjson.GetBytes(sanitized, "tools.0.parameters.properties.config.examples.0.pattern").String())
	require.Equal(t, "(?=keep-const)", gjson.GetBytes(sanitized, "tools.0.parameters.properties.config.const.pattern").String())
	require.Equal(t, "(?=keep-enum)", gjson.GetBytes(sanitized, "tools.0.parameters.properties.config.enum.0.pattern").String())
placeholder

func buildToolSchemaNullTypeBody(t *testing.T, hits int) []byte {
placeholder
	tools := make([]any, 0, hits)
	for i := 0; i < hits; i++ {
		tools = append(tools, map[string]any{
			"type":       "function",
			"name":       "automation_update",
			"parameters": map[string]any{"type": nil, "properties": map[string]any{placeholderplaceholder,
	placeholder)
placeholder
	body, err := json.Marshal(map[string]any{"model": "gpt-5.6-sol", "tools": toolsplaceholder)
placeholder
	return body
placeholder

// 复杂度守卫：重写次数必须与命中数无关。
//
// 逐个 sjson.SetBytes 的写法每命中一处就重扫并全量拷贝一次文档，命中 N 处即 N 次
// 全量拷贝；/v1/responses 的 body 上限是 gateway.max_body_size（默认 256MB），
// 构造请求可以塞进百万级命中，会被放大成 TB 级 memcpy。这里用分配次数锁死该行为：
// 命中数放大 500 倍，分配次数不得随之增长。
func TestSanitizeOpenAIResponsesToolParameterTypes_RewriteCountIndependentOfHits(t *testing.T) {
	small := buildToolSchemaNullTypeBody(t, 4)
	large := buildToolSchemaNullTypeBody(t, 2000)

	smallAllocs := testing.AllocsPerRun(2, func() {
		_, _, _ = sanitizeOpenAIResponsesToolParameterTypes(small)
placeholder)
	largeAllocs := testing.AllocsPerRun(2, func() {
		_, _, _ = sanitizeOpenAIResponsesToolParameterTypes(large)
placeholder)

	// 命中切片扩容是对数级，留出充裕余量；线性写法在这里会是 2000 量级。
	// 干净环境实测 large 约 17 allocs，200 是 10 倍余量，同时容忍 CI 慢 pod 上
	// 包内后台 goroutine（日志/ticker）对进程级 Mallocs 的噪声污染。
	require.Less(t, largeAllocs, 200.0,
		"分配次数随命中数线性增长，说明退回了逐路径全量重写 (small=%v large=%v)", smallAllocs, largeAllocs)

	// 同时确认大 body 的结果确实全部修好了。
	sanitized, changed, err := sanitizeOpenAIResponsesToolParameterTypes(large)
placeholder
	require.True(t, changed)
	require.Equal(t, 2000, int(gjson.GetBytes(sanitized, "tools.#").Int()))
	gjson.GetBytes(sanitized, "tools").ForEach(func(_, tool gjson.Result) bool {
		require.Equal(t, "object", tool.Get("parameters.type").String())
		return true
placeholder)
placeholder

func BenchmarkSanitizeOpenAIResponsesToolSchemaPatterns_ByteSpanPatch(b *testing.B) {
	var tools []any
	for i := 0; i < 2000; i++ {
		tools = append(tools, map[string]any{
			"type": "function",
			"name": "search",
			"parameters": map[string]any{
				"type":       "object",
				"properties": map[string]any{"query": map[string]any{"type": "string", "pattern": "(?=foo).*"placeholderplaceholder,
		placeholder,
	placeholder)
placeholder
	body, _ := json.Marshal(map[string]any{"model": "gpt-5.6", "tools": toolsplaceholder)
	b.ReportAllocs()
	b.SetBytes(int64(len(body)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = sanitizeOpenAIResponsesToolSchemaPatterns(body)
placeholder
placeholder

func BenchmarkSanitizeOpenAIResponsesToolParameterTypes_ByteSpanPatch(b *testing.B) {
	body := make([]byte, 0, 256*2000)
	body = append(body, []byte(`{"tools":[`)...)
	for i := 0; i < 2000; i++ {
		if i > 0 {
			body = append(body, ',')
	placeholder
		body = append(body, []byte(`{"type":"function","name":"automation_update","parameters":{"properties":{placeholder,"type":nullplaceholderplaceholder`)...)
placeholder
	body = append(body, []byte(`]placeholder`)...)
	b.ReportAllocs()
	b.SetBytes(int64(len(body)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = sanitizeOpenAIResponsesToolParameterTypes(body)
placeholder
placeholder
