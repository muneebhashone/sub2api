package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIsOpenAIWSTokenEvent_TerminalEventsExcluded 覆盖 isOpenAIWSTokenEvent 的回归用例。
// 重点验证终止事件（response.completed / response.done）不再被当作 token event，
// 否则当上游没有可识别的 delta 时，firstTokenMs 会被填到终止时刻，
// 等于把"总耗时"误报为"首 token 延迟"（issue #2651）。
func TestIsOpenAIWSTokenEvent_TerminalEventsExcluded(t *testing.T) {
	cases := []struct {
		name      string
		eventType string
		want      bool
placeholder{
		{name: "empty", eventType: "", want: falseplaceholder,
		{name: "whitespace_trimmed_empty", eventType: "   ", want: falseplaceholder,

		{name: "response.created", eventType: "response.created", want: falseplaceholder,
		{name: "response.in_progress", eventType: "response.in_progress", want: falseplaceholder,
		{name: "response.output_item.added", eventType: "response.output_item.added", want: falseplaceholder,
		{name: "response.output_item.done", eventType: "response.output_item.done", want: falseplaceholder,

		{name: "terminal_response.completed", eventType: "response.completed", want: falseplaceholder,
		{name: "terminal_response.done", eventType: "response.done", want: falseplaceholder,
		{name: "terminal_response.completed_padded", eventType: "  response.completed  ", want: falseplaceholder,
		{name: "terminal_response.done_padded", eventType: "  response.done  ", want: falseplaceholder,

		{name: "delta_text", eventType: "response.output_text.delta", want: trueplaceholder,
		{name: "delta_audio_transcript", eventType: "response.audio_transcript.delta", want: trueplaceholder,
		{name: "delta_function_call_arguments", eventType: "response.function_call_arguments.delta", want: trueplaceholder,

		{name: "output_text_done", eventType: "response.output_text.done", want: trueplaceholder,
		{name: "output_text_annotation_added", eventType: "response.output_text.annotation.added", want: trueplaceholder,

		{name: "output_audio_done", eventType: "response.output_audio.done", want: trueplaceholder,

		{name: "reasoning_summary_delta", eventType: "response.reasoning_summary_text.delta", want: trueplaceholder,

		{name: "unrelated_event_error", eventType: "error", want: falseplaceholder,
		{name: "unknown_event_without_match", eventType: "response.reasoning_summary_part.added", want: falseplaceholder,
placeholder

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := isOpenAIWSTokenEvent(tc.eventType)
			require.Equal(t, tc.want, got, "isOpenAIWSTokenEvent(%q)", tc.eventType)
	placeholder)
placeholder
placeholder

// TestIsOpenAIWSTokenEvent_DisjointWithTerminal 守护「token 事件集合与终止事件集合互斥」的不变量。
// firstTokenMs 的计算依赖于 isTokenEvent && !isTerminalEvent；
// 若两者再次出现交集，则 issue #2651 描述的 latency 误报会重现。
func TestIsOpenAIWSTokenEvent_DisjointWithTerminal(t *testing.T) {
	terminalEvents := []string{
		"response.completed",
		"response.done",
		"response.failed",
		"response.incomplete",
		"response.cancelled",
		"response.canceled",
placeholder
	for _, ev := range terminalEvents {
		ev := ev
		t.Run(ev, func(t *testing.T) {
			require.True(t, isOpenAIWSTerminalEvent(ev), "expected terminal event %q to be classified as terminal", ev)
			require.False(t, isOpenAIWSTokenEvent(ev), "terminal event %q must NOT be classified as token event (issue #2651)", ev)
	placeholder)
placeholder
placeholder
