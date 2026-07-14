package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsImageGenerationIntent(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		model    string
		body     []byte
		want     bool
placeholder{
		{
			name:     "images endpoint",
			endpoint: "/v1/images/generations",
			body:     []byte(`{"model":"gpt-image-2"placeholder`),
			want:     true,
	placeholder,
		{
			name:     "image model",
			endpoint: "/v1/responses",
			model:    "gpt-image-2",
			body:     []byte(`{"model":"gpt-image-2"placeholder`),
			want:     true,
	placeholder,
		{
			name:     "image tool",
			endpoint: "/v1/responses",
			model:    "gpt-5.4",
			body:     []byte(`{"model":"gpt-5.4","tools":[{"type":"image_generation"placeholder]placeholder`),
			want:     true,
	placeholder,
		{
			name:     "image tool choice",
			endpoint: "/v1/responses",
			model:    "gpt-5.4",
			body:     []byte(`{"model":"gpt-5.4","tool_choice":{"type":"image_generation"placeholderplaceholder`),
			want:     true,
	placeholder,
		{
			name:     "namespace image_gen tool choice",
			endpoint: "/v1/responses",
			model:    "gpt-5.5",
			body:     []byte(`{"model":"gpt-5.5","tool_choice":{"type":"namespace","name":"image_gen"placeholderplaceholder`),
			want:     true,
	placeholder,
		{
			name:     "custom imagegen function tool choice is not image intent",
			endpoint: "/v1/responses",
			model:    "gpt-5.5",
			body:     []byte(`{"model":"gpt-5.5","tool_choice":{"function":{"name":"imagegen"placeholderplaceholderplaceholder`),
			want:     false,
	placeholder,
		{
			name:     "required tool choice alone is text",
			endpoint: "/v1/responses",
			model:    "gpt-5.4",
			body:     []byte(`{"model":"gpt-5.4","tool_choice":"required"placeholder`),
			want:     false,
	placeholder,
		{
			name:     "text only gpt 5.4",
			endpoint: "/v1/responses",
			model:    "gpt-5.4",
			body:     []byte(`{"model":"gpt-5.4","input":"write code"placeholder`),
			want:     false,
	placeholder,
		{
			name:     "namespace image_gen tool in top-level tools",
			endpoint: "/v1/responses",
			model:    "gpt-5.5",
			body:     []byte(`{"model":"gpt-5.5","tools":[{"type":"namespace","name":"image_gen","tools":[{"type":"function","name":"imagegen"placeholder]placeholder]placeholder`),
			want:     true,
	placeholder,
		{
			name:     "custom namespace with nested imagegen function is not image intent",
			endpoint: "/v1/responses",
			model:    "gpt-5.5",
			body:     []byte(`{"model":"gpt-5.5","tools":[{"type":"namespace","name":"media_tools","tools":[{"type":"function","name":"imagegen"placeholder]placeholder]placeholder`),
			want:     false,
	placeholder,
		{
			name:     "namespace image_gen in input additional_tools (Responses Lite)",
			endpoint: "/v1/responses",
			model:    "gpt-5.5",
			body:     []byte(`{"model":"gpt-5.5","input":[{"type":"additional_tools","role":"developer","tools":[{"type":"namespace","name":"image_gen","tools":[{"type":"function","name":"imagegen"placeholder]placeholder]placeholder]placeholder`),
			want:     true,
	placeholder,
		{
			name:     "non-image namespace tool is not flagged",
			endpoint: "/v1/responses",
			model:    "gpt-5.5",
			body:     []byte(`{"model":"gpt-5.5","tools":[{"type":"namespace","name":"code_tools","tools":[{"type":"function","name":"run"placeholder]placeholder]placeholder`),
			want:     false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsImageGenerationIntent(tt.endpoint, tt.model, tt.body))
	placeholder)
placeholder
placeholder

func TestIsImageGenerationIntentJSONSemantics(t *testing.T) {
	largeInput := strings.Repeat("x", 1<<20)
	tests := []struct {
		name     string
		endpoint string
		body     []byte
		want     bool
placeholder{
		{
			name:     "chat body image model",
			endpoint: "/v1/chat/completions",
			body:     []byte(`{"model":"gpt-image-2"placeholder`),
			want:     true,
	placeholder,
		{
			name:     "large responses input with trailing namespace tool choice",
			endpoint: "/v1/responses",
			body:     []byte(`{"model":"gpt-5.5","input":"` + largeInput + `","tool_choice":{"type":"namespace","name":"image_gen"placeholderplaceholder`),
			want:     true,
	placeholder,
		{
			name:     "invalid json with image tool",
			endpoint: "/v1/responses",
			body:     []byte(`{"tools":[{"type":"image_generation"placeholder]`),
			want:     false,
	placeholder,
		{
			name:     "duplicate model uses first value",
			endpoint: "/v1/responses",
			body:     []byte(`{"model":"gpt-5.5","model":"gpt-image-2"placeholder`),
			want:     false,
	placeholder,
		{
			name:     "duplicate null model still uses first value",
			endpoint: "/v1/responses",
			body:     []byte(`{"model":null,"model":"gpt-image-2"placeholder`),
			want:     false,
	placeholder,
		{
			name:     "duplicate tools uses first value",
			endpoint: "/v1/responses",
			body:     []byte(`{"tools":[],"tools":[{"type":"image_generation"placeholder]placeholder`),
			want:     false,
	placeholder,
		{
			name:     "duplicate input uses first value",
			endpoint: "/v1/responses",
			body:     []byte(`{"input":[],"input":[{"type":"additional_tools","tools":[{"type":"namespace","name":"image_gen"placeholder]placeholder]placeholder`),
			want:     false,
	placeholder,
		{
			name:     "duplicate tool choice uses first value",
			endpoint: "/v1/responses",
			body:     []byte(`{"tool_choice":"required","tool_choice":{"type":"image_generation"placeholderplaceholder`),
			want:     false,
	placeholder,
		{
			name:     "escaped top level key",
			endpoint: "/v1/responses",
			body:     []byte(`{"tool_\u0063hoice":{"type":"image_generation"placeholderplaceholder`),
			want:     true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsImageGenerationIntent(tt.endpoint, "gpt-5.5", tt.body))
	placeholder)
placeholder
placeholder

func TestIsImageGenerationIntentMap_NamespaceImageGen(t *testing.T) {
	tests := []struct {
		name    string
		reqBody map[string]any
		want    bool
placeholder{
		{
			name: "top-level namespace image_gen",
			reqBody: map[string]any{
				"model": "gpt-5.5",
				"tools": []any{
					map[string]any{"type": "namespace", "name": "image_gen", "tools": []any{
						map[string]any{"type": "function", "name": "imagegen"placeholder,
			placeholder
			placeholder,
		placeholder,
			want: true,
	placeholder,
		{
			name: "additional_tools in input",
			reqBody: map[string]any{
				"model": "gpt-5.5",
				"input": []any{
					map[string]any{
						"type": "additional_tools",
						"tools": []any{
							map[string]any{"type": "namespace", "name": "image_gen"placeholder,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			want: true,
	placeholder,
		{
			name: "custom namespace with nested imagegen function is not image intent",
			reqBody: map[string]any{
				"model": "gpt-5.5",
				"tools": []any{
					map[string]any{
						"type": "namespace",
						"name": "media_tools",
						"tools": []any{
							map[string]any{"type": "function", "name": "imagegen"placeholder,
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			want: false,
	placeholder,
		{
			name: "namespace image_gen tool choice",
			reqBody: map[string]any{
				"model":       "gpt-5.5",
				"tool_choice": map[string]any{"type": "namespace", "name": "image_gen"placeholder,
		placeholder,
			want: true,
	placeholder,
		{
			name: "custom imagegen function tool choice is not image intent",
			reqBody: map[string]any{
				"model": "gpt-5.5",
				"tool_choice": map[string]any{
					"function": map[string]any{"name": "imagegen"placeholder,
			placeholder,
		placeholder,
			want: false,
	placeholder,
		{
			name: "non-image namespace not flagged",
			reqBody: map[string]any{
				"model": "gpt-5.5",
				"tools": []any{
					map[string]any{"type": "namespace", "name": "code_tools"placeholder,
			placeholder,
		placeholder,
			want: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsImageGenerationIntentMap("/v1/responses", "gpt-5.5", tt.reqBody))
	placeholder)
placeholder
placeholder

func TestResolveOpenAIResponsesImageBillingConfigUsesCurrentBodyModel(t *testing.T) {
	imageModel, imageSize, err := resolveOpenAIResponsesImageBillingConfigFromBody(
		[]byte(`{"model":"mapped-image-model","tools":[{"type":"image_generation","size":"1024x1024"placeholder]placeholder`),
		"requested-model",
	)
placeholder
	require.Equal(t, "mapped-image-model", imageModel)
	require.Equal(t, "1K", imageSize)
placeholder

func TestResolveOpenAIResponsesImageBillingConfigToolModelWins(t *testing.T) {
	imageModel, imageSize, err := resolveOpenAIResponsesImageBillingConfigFromBody(
		[]byte(`{"model":"mapped-text-model","tools":[{"type":"image_generation","model":"gpt-image-2","size":"1536x1024"placeholder]placeholder`),
		"requested-model",
	)
placeholder
	require.Equal(t, "gpt-image-2", imageModel)
	require.Equal(t, "2K", imageSize)
placeholder

func TestResolveOpenAIResponsesImageBillingConfigFromBodyIgnoresUnrelatedLargeInput(t *testing.T) {
	cfg, err := resolveOpenAIResponsesImageBillingConfigDetailedFromBody(
		[]byte(`{"model":"mapped-text-model","tools":[{"type":"image_generation","model":"gpt-image-2","size":"2048x1152"placeholder],"input":[{"type":"message","content":[{"type":"input_text","text":"hi","nonce":1e1000000placeholder]placeholder]placeholder`),
		"requested-model",
	)
placeholder
	require.Equal(t, "gpt-image-2", cfg.Model)
	require.Equal(t, "2K", cfg.SizeTier)
	require.Equal(t, "2048x1152", cfg.InputSize)
placeholder

func TestResolveOpenAIResponsesImageBillingConfigSupportsOfficialAndCustomSizes(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		wantTier string
placeholder{
		{
			name:     "official 2k landscape",
			body:     []byte(`{"model":"gpt-5.4","tools":[{"type":"image_generation","model":"gpt-image-2","size":"2048x1152"placeholder]placeholder`),
			wantTier: "2K",
	placeholder,
		{
			name:     "official 4k landscape",
			body:     []byte(`{"model":"gpt-5.4","tools":[{"type":"image_generation","model":"gpt-image-2","size":"3840x2160"placeholder]placeholder`),
			wantTier: "4K",
	placeholder,
		{
			name:     "custom valid 2k",
			body:     []byte(`{"model":"gpt-5.5","tools":[{"type":"image_generation","model":"gpt-image-2","size":"1280x768"placeholder]placeholder`),
			wantTier: "2K",
	placeholder,
		{
			name:     "default image tool model supports flexible size",
			body:     []byte(`{"model":"gpt-5.4","tools":[{"type":"image_generation","size":"2048x1152"placeholder]placeholder`),
			wantTier: "2K",
	placeholder,
		{
			name:     "top level image size is moved into billing",
			body:     []byte(`{"model":"gpt-image-2","size":"2048x2048","tools":[{"type":"image_generation","model":"gpt-image-2"placeholder]placeholder`),
			wantTier: "2K",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageModel, imageSize, err := resolveOpenAIResponsesImageBillingConfigFromBody(tt.body, "requested-model")
		placeholder
			require.NotEmpty(t, imageModel)
			require.Equal(t, tt.wantTier, imageSize)
	placeholder)
placeholder
placeholder

func TestResolveOpenAIResponsesImageBillingConfigDoesNotRejectUnknownSizes(t *testing.T) {
	imageModel, imageSize, err := resolveOpenAIResponsesImageBillingConfigFromBody(
		[]byte(`{"model":"gpt-5.4","tools":[{"type":"image_generation","model":"gpt-image-1.5","size":"2048x1152"placeholder]placeholder`),
		"requested-model",
	)
placeholder
	require.Equal(t, "gpt-image-1.5", imageModel)
	require.Equal(t, "2K", imageSize)
placeholder

func TestOpenAIImageOutputCounterDeduplicatesFinalImages(t *testing.T) {
	counter := newOpenAIImageOutputCounter()
	counter.AddSSEData([]byte(`{"type":"response.image_generation_call.partial_image","partial_image_b64":"abc"placeholder`))
	counter.AddSSEData([]byte(`{"type":"response.output_item.done","item":{"id":"ig_1","type":"image_generation_call","result":"final-a","size":"1024x1024"placeholderplaceholder`))
	counter.AddSSEData([]byte(`{"type":"response.completed","response":{"output":[{"id":"ig_1","type":"image_generation_call","result":"final-a"placeholder,{"id":"ig_2","type":"image_generation_call","result":"final-b","size":"3840x2160"placeholder]placeholderplaceholder`))
	require.Equal(t, 2, counter.Count())
	require.Equal(t, []string{"1024x1024", "3840x2160"placeholder, counter.Sizes())
placeholder

func TestOpenAIImageOutputCounterCountsImagesAPIStreamShapes(t *testing.T) {
	counter := newOpenAIImageOutputCounter()
	counter.AddSSEData([]byte(`{"type":"image_generation.completed","id":"ig_complete","b64_json":"final-a"placeholder`))
	counter.AddSSEData([]byte(`{"type":"response.output_item.done","item":{"id":"ig_item","type":"image_generation_call","result":"final-b"placeholderplaceholder`))
	counter.AddSSEData([]byte(`{"type":"response.completed","response":{"output":[{"id":"ig_done","type":"image_generation_call","result":"final-c"placeholder]placeholderplaceholder`))
	require.Equal(t, 3, counter.Count())

	dataCounter := newOpenAIImageOutputCounter()
	dataCounter.AddSSEData([]byte(`{"data":[{"b64_json":"a"placeholder,{"b64_json":"b"placeholder]placeholder`))
	dataCounter.AddSSEData([]byte(`{"data":[{"b64_json":"a"placeholder,{"b64_json":"b"placeholder,{"b64_json":"c"placeholder]placeholder`))
	require.Equal(t, 3, dataCounter.Count())
placeholder

func TestOpenAIImageOutputCounterCountsMultilineSSEDataPayload(t *testing.T) {
	counter := newOpenAIImageOutputCounter()
	counter.AddSSEData([]byte("{\"type\":\"image_generation.completed\",\n\"b64_json\":\"final-a\"placeholder"))
	require.Equal(t, 1, counter.Count())
placeholder

func TestOpenAIImageOutputCounterCountsMultilineSSEBodyPayload(t *testing.T) {
	counter := newOpenAIImageOutputCounter()
	counter.AddSSEBody(
		"data: {\"type\":\"image_generation.completed\",\n" +
			"data: \"b64_json\":\"final-a\"placeholder\n\n" +
			"data: [DONE]\n\n",
	)
	require.Equal(t, 1, counter.Count())
placeholder

func TestOpenAIImageOutputCounterFallsBackForInvalidMultilineSSEBody(t *testing.T) {
	counter := newOpenAIImageOutputCounter()
	counter.AddSSEBody(
		"data: {\"type\":\"image_generation.completed\",\"b64_json\":\"final-a\"placeholder\n" +
			"data: {\"type\":\"image_generation.completed\",\"b64_json\":\"final-b\"placeholder\n\n",
	)
	require.Equal(t, 2, counter.Count())
placeholder

func TestCollectOpenAIResponseImageOutputSizesFromJSONBytes(t *testing.T) {
	body := []byte(`{
		"output": [
			{"id":"ig_1","type":"image_generation_call","result":"final-a","size":"3840x2160"placeholder,
			{"id":"ig_2","type":"image_generation_call","result":"final-b","size":"1024x1024"placeholder
		]
placeholder`)

	require.Equal(t, 2, countOpenAIResponseImageOutputsFromJSONBytes(body))
	require.Equal(t, []string{"3840x2160", "1024x1024"placeholder, collectOpenAIResponseImageOutputSizesFromJSONBytes(body))
placeholder

func TestCollectOpenAIResponseImageOutputSizesFromImagesAPIData(t *testing.T) {
	body := []byte(`{
		"data": [
			{"b64_json":"final-a","size":"2048x1152"placeholder,
			{"b64_json":"final-b","size":"2048x1152"placeholder
		]
placeholder`)

	require.Equal(t, 2, countOpenAIResponseImageOutputsFromJSONBytes(body))
	require.Equal(t, []string{"2048x1152", "2048x1152"placeholder, collectOpenAIResponseImageOutputSizesFromJSONBytes(body))
placeholder

func TestCollectOpenAIImageOutputSizesFromSSEBody(t *testing.T) {
	body := "data: {\"type\":\"response.output_item.done\",\"item\":{\"id\":\"ig_1\",\"type\":\"image_generation_call\",\"result\":\"final-a\",\"size\":\"3840x2160\"placeholderplaceholder\n\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"output\":[{\"id\":\"ig_1\",\"type\":\"image_generation_call\",\"result\":\"final-a\"placeholder,{\"id\":\"ig_2\",\"type\":\"image_generation_call\",\"result\":\"final-b\",\"size\":\"1024x1024\"placeholder]placeholderplaceholder\n\n" +
		"data: [DONE]\n\n"

	require.Equal(t, 2, countOpenAIImageOutputsFromSSEBody(body))
	require.Equal(t, []string{"3840x2160", "1024x1024"placeholder, collectOpenAIImageOutputSizesFromSSEBody(body))
placeholder
