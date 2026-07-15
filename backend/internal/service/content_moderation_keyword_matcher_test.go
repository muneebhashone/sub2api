package service

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContentModerationKeywordMatcherMatchesLegacyBehavior(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		keywords []string
placeholder{
		{name: "miss", text: "clean prompt", keywords: []string{"blocked", "secret"placeholderplaceholder,
		{name: "case insensitive", text: "contains SECRET value", keywords: []string{"secret"placeholderplaceholder,
		{name: "configured order wins", text: "early appears before later", keywords: []string{"later", "early"placeholderplaceholder,
		{name: "overlap uses configured order", text: "abc", keywords: []string{"bc", "abc"placeholderplaceholder,
		{name: "unicode", text: "这里包含敏感词和世界", keywords: []string{"世界", "敏感词"placeholderplaceholder,
		{name: "duplicates", text: "duplicate", keywords: []string{"duplicate", "DUPLICATE"placeholderplaceholder,
		{name: "empty entries", text: "blocked", keywords: []string{"", "blocked"placeholderplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantKeyword, wantHit := matchBlockedKeyword(tt.text, tt.keywords)
			gotKeyword, gotHit := newContentModerationKeywordMatcher(tt.keywords).Match(tt.text)
			require.Equal(t, wantHit, gotHit)
			require.Equal(t, wantKeyword, gotKeyword)
	placeholder)
placeholder
placeholder

func TestContentModerationKeywordMatcherRandomizedParity(t *testing.T) {
	rng := rand.New(rand.NewSource(20260714))
	const alphabet = "abcXYZ"
	for iteration := 0; iteration < 1000; iteration++ {
		keywords := make([]string, 1+rng.Intn(30))
		for index := range keywords {
			length := 1 + rng.Intn(8)
			var value strings.Builder
			for range length {
				_ = value.WriteByte(alphabet[rng.Intn(len(alphabet))])
		placeholder
			keywords[index] = value.String()
	placeholder
		var text strings.Builder
		for range 20 + rng.Intn(100) {
			_ = text.WriteByte(alphabet[rng.Intn(len(alphabet))])
	placeholder

		wantKeyword, wantHit := matchBlockedKeyword(text.String(), keywords)
		gotKeyword, gotHit := newContentModerationKeywordMatcher(keywords).Match(text.String())
		require.Equal(t, wantHit, gotHit, "iteration %d", iteration)
		require.Equal(t, wantKeyword, gotKeyword, "iteration %d", iteration)
placeholder
placeholder
