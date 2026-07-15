package service

import (
	"strings"
	"testing"
)

var benchmarkOpenAIContentSessionSeed string

func BenchmarkDeriveOpenAIContentSessionSeedLargeBody(b *testing.B) {
	largeHistory := strings.Repeat("payload", 1<<17)
	tests := []struct {
		name string
		body []byte
placeholder{
		{name: "ChatCompletions", body: []byte(`{"model":"gpt-5.4","tools":[{"type":"function","function":{"name":"lookup"placeholderplaceholder],"messages":[{"role":"system","content":"Be concise."placeholder,{"role":"user","content":"Hello"placeholder,{"role":"assistant","content":"` + largeHistory + `"placeholder,{"role":"user","content":"Follow-up"placeholder]placeholder`)placeholder,
		{name: "Responses", body: []byte(`{"model":"gpt-5.4","instructions":"Be concise.","tools":[{"type":"function","name":"lookup"placeholder],"input":[{"role":"system","content":"System prompt"placeholder,{"role":"user","content":"Hello"placeholder,{"role":"assistant","content":"` + largeHistory + `"placeholder]placeholder`)placeholder,
placeholder

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(test.body)))
			for range b.N {
				benchmarkOpenAIContentSessionSeed = deriveOpenAIContentSessionSeed(test.body)
		placeholder
	placeholder)
placeholder
placeholder
