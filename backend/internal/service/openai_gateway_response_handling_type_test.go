package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestOpenAIStreamEventIsTerminalWithTypeMatchesExistingSemantics(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
placeholder{
		{name: "empty", data: "", want: falseplaceholder,
		{name: "whitespace", data: " \t ", want: falseplaceholder,
		{name: "done", data: " [DONE] ", want: trueplaceholder,
		{name: "JSON outer whitespace", data: " \n\t {\"type\":\"response.completed\"placeholder \r\n", want: trueplaceholder,
		{name: "completed", data: `{"type":"response.completed"placeholder`, want: trueplaceholder,
		{name: "response done", data: `{"type":"response.done"placeholder`, want: trueplaceholder,
		{name: "failed", data: `{"type":"response.failed"placeholder`, want: trueplaceholder,
		{name: "incomplete", data: `{"type":"response.incomplete"placeholder`, want: trueplaceholder,
		{name: "cancelled", data: `{"type":"response.cancelled"placeholder`, want: trueplaceholder,
		{name: "canceled", data: `{"type":"response.canceled"placeholder`, want: trueplaceholder,
		{name: "delta", data: `{"type":"response.output_text.delta"placeholder`, want: falseplaceholder,
		{name: "invalid JSON", data: `{"type":`, want: falseplaceholder,
		{name: "terminal with trailing garbage", data: `{"type":"response.completed"placeholder trailing`, want: trueplaceholder,
		{name: "nonterminal with trailing garbage", data: `{"type":"response.output_text.delta"placeholder trailing`, want: falseplaceholder,
		{name: "type whitespace remains nonterminal", data: `{"type":" response.completed "placeholder`, want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventType := gjson.GetBytes([]byte(tt.data), "type").String()
			got := openAIStreamEventIsTerminalWithType(tt.data, eventType)

			require.Equal(t, tt.want, got)
			require.Equal(t, openAIStreamEventIsTerminal(tt.data), got)
	placeholder)
placeholder
placeholder

var (
	benchmarkOpenAIResponseSSEEventTypeSink string
	benchmarkOpenAIResponseSSETerminalSink  bool
)

func BenchmarkOpenAIResponseSSETypeExtraction(b *testing.B) {
	data := `{"type":"response.output_text.delta","sequence_number":42,"delta":"streaming response benchmark payload"placeholder`
	dataBytes := []byte(data)

	b.Run("legacy double parse", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(dataBytes)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkOpenAIResponseSSETerminalSink = openAIStreamEventIsTerminal(data)
			benchmarkOpenAIResponseSSEEventTypeSink = strings.TrimSpace(gjson.GetBytes(dataBytes, "type").String())
	placeholder
placeholder)

	b.Run("reused single parse", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(dataBytes)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			eventTypeRaw := gjson.GetBytes(dataBytes, "type").String()
			benchmarkOpenAIResponseSSEEventTypeSink = strings.TrimSpace(eventTypeRaw)
			benchmarkOpenAIResponseSSETerminalSink = openAIStreamEventIsTerminalWithType(data, eventTypeRaw)
	placeholder
placeholder)
placeholder
