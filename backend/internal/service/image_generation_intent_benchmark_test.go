package service

import (
	"strings"
	"testing"
)

var imageGenerationIntentBenchmarkResult bool

func BenchmarkIsImageGenerationIntent(b *testing.B) {
	largeInput := strings.Repeat("x", 1<<20)
	benchmarks := []struct {
		name string
		body []byte
		want bool
placeholder{
		{
			name: "1MiBInputNoImage",
			body: []byte(`{"model":"gpt-5.5","tools":[],"input":"` + largeInput + `","tool_choice":"auto"placeholder`),
			want: false,
	placeholder,
		{
			name: "1MiBInputLeadingImageTool",
			body: []byte(`{"model":"gpt-5.5","tools":[{"type":"image_generation"placeholder],"input":"` + largeInput + `"placeholder`),
			want: true,
	placeholder,
		{
			name: "1MiBInputTrailingToolChoice",
			body: []byte(`{"model":"gpt-5.5","tools":[],"input":"` + largeInput + `","tool_choice":{"type":"image_generation"placeholderplaceholder`),
			want: true,
	placeholder,
		{
			name: "Invalid1MiBJSON",
			body: []byte(`{"model":"gpt-5.5","input":"` + largeInput),
			want: false,
	placeholder,
		{
			name: "DuplicateKeysFirstWins",
			body: []byte(`{"model":"gpt-5.5","model":"gpt-image-2","tools":[],"tools":[{"type":"image_generation"placeholder],"input":"` + largeInput + `","tool_choice":"auto","tool_choice":{"type":"image_generation"placeholderplaceholder`),
			want: false,
	placeholder,
placeholder

	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			if got := IsImageGenerationIntent("/v1/responses", "gpt-5.5", benchmark.body); got != benchmark.want {
				b.Fatalf("IsImageGenerationIntent() = %v, want %v", got, benchmark.want)
		placeholder
			b.ReportAllocs()
			b.SetBytes(int64(len(benchmark.body)))
			b.ResetTimer()
			var result bool
			for i := 0; i < b.N; i++ {
				result = IsImageGenerationIntent("/v1/responses", "gpt-5.5", benchmark.body)
		placeholder
			imageGenerationIntentBenchmarkResult = result
	placeholder)
placeholder
placeholder
