package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestNormalizeMaxReasoningEffort(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
placeholder{
		{name: "empty", in: "", want: ""placeholder,
		{name: "separator", in: "x-high", want: "xhigh"placeholder,
		{name: "max is distinct", in: "max", want: "max"placeholder,
		{name: "none is unsupported", in: "none", want: ""placeholder,
		{name: "invalid", in: "banana", want: ""placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NormalizeMaxReasoningEffort(tt.in))
	placeholder)
placeholder
placeholder

func TestNormalizeReasoningEffortMappings(t *testing.T) {
	t.Run("canonicalizes fixed OpenAI values", func(t *testing.T) {
		got, err := NormalizeReasoningEffortMappings(PlatformOpenAI, []ReasoningEffortMapping{
			{From: " MAX ", To: " x-high "placeholder,
			{From: "minimal", To: "high"placeholder,
	placeholder)
	placeholder
		require.Equal(t, []ReasoningEffortMapping{
			{From: "max", To: "xhigh"placeholder,
			{From: "minimal", To: "high"placeholder,
	placeholder, got)
placeholder)

	t.Run("rejects empty values", func(t *testing.T) {
		_, err := NormalizeReasoningEffortMappings(PlatformOpenAI, []ReasoningEffortMapping{{From: "max"placeholderplaceholder)
		require.ErrorContains(t, err, "empty or unknown")
placeholder)

	t.Run("rejects duplicate sources case insensitively", func(t *testing.T) {
		_, err := NormalizeReasoningEffortMappings(PlatformOpenAI, []ReasoningEffortMapping{
			{From: "max", To: "xhigh"placeholder,
			{From: " MAX ", To: "high"placeholder,
	placeholder)
		require.ErrorContains(t, err, "duplicate")
placeholder)

	t.Run("rejects mappings for non OpenAI platforms", func(t *testing.T) {
		for _, platform := range []string{PlatformAnthropic, PlatformGemini, PlatformAntigravity, PlatformGrokplaceholder {
			_, err := NormalizeReasoningEffortMappings(platform, []ReasoningEffortMapping{{From: "low", To: "high"placeholderplaceholder)
			require.ErrorContains(t, err, "only supported for platform \"openai\"")
	placeholder

		_, err := NormalizeReasoningEffortMappings(PlatformOpenAI, []ReasoningEffortMapping{{From: "none", To: "low"placeholderplaceholder)
		require.ErrorContains(t, err, "empty or unknown")

		_, err = NormalizeReasoningEffortMappings(PlatformOpenAI, []ReasoningEffortMapping{{From: "ultra", To: "high"placeholderplaceholder)
		require.ErrorContains(t, err, "empty or unknown")
placeholder)
placeholder

func TestNormalizeMaxReasoningEffortForPlatform(t *testing.T) {
	value, err := normalizeMaxReasoningEffortForPlatform(PlatformOpenAI, "max")
placeholder
	require.Equal(t, "max", value)

	for _, platform := range []string{PlatformAnthropic, PlatformGemini, PlatformAntigravity, PlatformGrokplaceholder {
		_, err = normalizeMaxReasoningEffortForPlatform(platform, "low")
		require.ErrorContains(t, err, "only supported for platform \"openai\"")
placeholder

	_, err = normalizeMaxReasoningEffortForPlatform(PlatformOpenAI, "none")
	require.ErrorContains(t, err, "not supported")
placeholder

func TestApplyOpenAIReasoningEffortPolicy(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		max      string
		mappings []ReasoningEffortMapping
		path     string
		want     string
		changed  bool
placeholder{
		{name: "nested caps high", body: `{"reasoning":{"effort":"xhigh"placeholderplaceholder`, max: "medium", path: "reasoning.effort", want: "medium", changed: trueplaceholder,
		{name: "flat caps high", body: `{"reasoning_effort":"high"placeholder`, max: "low", path: "reasoning_effort", want: "low", changed: trueplaceholder,
		{name: "does not raise omitted", body: `{"model":"gpt-5"placeholder`, max: "low", path: "reasoning_effort", want: "", changed: falseplaceholder,
		{name: "keeps lower value", body: `{"reasoning_effort":"low"placeholder`, max: "high", path: "reasoning_effort", want: "low", changed: falseplaceholder,
		{name: "normalizes request alias", body: `{"reasoning_effort":"x-high"placeholder`, max: "xhigh", path: "reasoning_effort", want: "xhigh", changed: trueplaceholder,
		{name: "caps max below its distinct rank", body: `{"reasoning_effort":"max"placeholder`, max: "xhigh", path: "reasoning_effort", want: "xhigh", changed: trueplaceholder,
		{name: "keeps xhigh below max", body: `{"reasoning_effort":"xhigh"placeholder`, max: "max", path: "reasoning_effort", want: "xhigh", changed: falseplaceholder,
		{name: "ignores stale none ceiling", body: `{"reasoning_effort":"high"placeholder`, max: "none", path: "reasoning_effort", want: "high", changed: falseplaceholder,
		{name: "caps both shapes", body: `{"reasoning":{"effort":"high"placeholder,"reasoning_effort":"xhigh"placeholder`, max: "low", path: "reasoning.effort", want: "low", changed: trueplaceholder,
		{name: "maps before cap", body: `{"reasoning":{"effort":"MAX"placeholderplaceholder`, max: "medium", mappings: []ReasoningEffortMapping{{From: "max", To: "xhigh"placeholderplaceholder, path: "reasoning.effort", want: "medium", changed: trueplaceholder,
		{name: "does not chain mappings", body: `{"reasoning_effort":"max"placeholder`, mappings: []ReasoningEffortMapping{{From: "max", To: "xhigh"placeholder, {From: "xhigh", To: "low"placeholderplaceholder, path: "reasoning_effort", want: "xhigh", changed: trueplaceholder,
		{name: "keeps unknown without mapping", body: `{"reasoning_effort":"future"placeholder`, max: "low", path: "reasoning_effort", want: "future", changed: falseplaceholder,
		{name: "keeps non string value", body: `{"reasoning_effort":{"level":"high"placeholderplaceholder`, max: "low", path: "reasoning_effort.level", want: "high", changed: falseplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, changed := ApplyOpenAIReasoningEffortPolicy([]byte(tt.body), tt.max, tt.mappings)
			require.Equal(t, tt.changed, changed)
			if tt.path != "" {
				require.Equal(t, tt.want, gjson.GetBytes(got, tt.path).String())
		placeholder
	placeholder)
placeholder
placeholder
