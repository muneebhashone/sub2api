package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsImageGenerationIntentForPlatform_GrokCodexDeclarations(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
placeholder{
		{
			name: "top-level image_gen namespace is passive",
			body: `{"model":"grok-4.5","tools":[{"type":"namespace","name":"image_gen","tools":[{"type":"function","name":"imagegen"placeholder]placeholder],"tool_choice":"auto","input":"write code"placeholder`,
	placeholder,
		{
			name: "responses lite additional_tools image_gen is passive",
			body: `{"model":"grok-4.5","tool_choice":"auto","input":[{"type":"additional_tools","tools":[{"type":"namespace","name":"image_gen","tools":[{"type":"function","name":"imagegen"placeholder]placeholder]placeholder,{"type":"message","role":"user","content":"write code"placeholder]placeholder`,
	placeholder,
		{
			name: "legacy adapter flattened image_gen function declaration is passive",
			body: `{"model":"grok-4.5","tools":[{"type":"function","name":"image_gen.imagegen"placeholder],"input":"write code"placeholder`,
	placeholder,
		{
			name: "native image generation declaration remains explicit",
			body: `{"model":"grok-4.5","tools":[{"type":"image_generation"placeholder],"input":"draw a cat"placeholder`,
			want: true,
	placeholder,
		{
			name: "native image generation tool choice is explicit",
			body: `{"model":"grok-4.5","tools":[{"type":"image_generation"placeholder],"tool_choice":{"type":"image_generation"placeholder,"input":"draw a cat"placeholder`,
			want: true,
	placeholder,
		{
			name: "image_gen namespace tool choice is explicit",
			body: `{"model":"grok-4.5","tools":[{"type":"namespace","name":"image_gen"placeholder],"tool_choice":{"type":"namespace","name":"image_gen"placeholder,"input":"draw a cat"placeholder`,
			want: true,
	placeholder,
		{
			name: "flattened image_gen function tool choice is explicit",
			body: `{"model":"grok-4.5","tools":[{"type":"function","name":"image_gen.imagegen"placeholder],"tool_choice":{"type":"function","name":"image_gen.imagegen"placeholder,"input":"draw a cat"placeholder`,
			want: true,
	placeholder,
		{
			name: "wrapped image_gen function tool choice is explicit",
			body: `{"model":"grok-4.5","tools":[{"type":"function","name":"image_gen.imagegen"placeholder],"tool_choice":{"tool":{"type":"function","name":"image_gen.imagegen"placeholderplaceholder,"input":"draw a cat"placeholder`,
			want: true,
	placeholder,
		{
			name: "vision input is not image generation",
			body: `{"model":"grok-4.5","input":[{"type":"message","role":"user","content":[{"type":"input_image","image_url":"data:image/png;base64,AA=="placeholder]placeholder]placeholder`,
	placeholder,
		{
			name: "ordinary function declaration is not image generation",
			body: `{"model":"grok-4.5","tools":[{"type":"function","name":"lookup"placeholder],"tool_choice":"auto","input":"write code"placeholder`,
	placeholder,
		{
			name: "image_gen function call history is not current intent",
			body: `{"model":"grok-4.5","input":[{"type":"function_call","namespace":"image_gen","name":"imagegen","arguments":"{placeholder"placeholder]placeholder`,
	placeholder,
		{
			name: "image generation call history is not current intent",
			body: `{"model":"grok-4.5","input":[{"type":"image_generation_call","id":"ig_1"placeholder]placeholder`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsImageGenerationIntentForPlatform(
				openAIResponsesEndpoint,
				"grok-4.5",
				[]byte(tt.body),
				PlatformGrok,
			))
	placeholder)
placeholder
placeholder

func TestIsImageGenerationIntentForPlatform_GrokPreservesHardSignals(t *testing.T) {
	require.True(t, IsImageGenerationIntentForPlatform(
		"/v1/images/generations",
		"grok-4.5",
		[]byte(`{"input":"draw"placeholder`),
		PlatformGrok,
	))
	require.True(t, IsImageGenerationIntentForPlatform(
		openAIResponsesEndpoint,
		"gpt-image-2",
		[]byte(`{"input":"draw"placeholder`),
		PlatformGrok,
	))
	require.True(t, IsImageGenerationIntentForPlatform(
		openAIResponsesEndpoint,
		"grok-4.5",
		[]byte(`{"model":"gpt-image-2","input":"draw"placeholder`),
		PlatformGrok,
	))
placeholder

func TestIsImageGenerationIntentForPlatform_OtherPlatformsKeepDeclarationSemantics(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","tools":[{"type":"namespace","name":"image_gen","tools":[{"type":"function","name":"imagegen"placeholder]placeholder],"input":"write code"placeholder`)

	require.True(t, IsImageGenerationIntentForPlatform(openAIResponsesEndpoint, "gpt-5.5", body, PlatformOpenAI))
	require.True(t, IsImageGenerationIntentForPlatform(openAIResponsesEndpoint, "gpt-5.5", body, PlatformAnthropic))
placeholder
