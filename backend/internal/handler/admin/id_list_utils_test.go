//go:build unit

package admin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeInt64IDList(t *testing.T) {
	tests := []struct {
		name string
		in   []int64
		want []int64
placeholder{
		{"nil input", nil, nilplaceholder,
		{"empty input", []int64{placeholder, nilplaceholder,
		{"single element", []int64{5placeholder, []int64{5placeholderplaceholder,
		{"already sorted unique", []int64{1, 2, 3placeholder, []int64{1, 2, 3placeholderplaceholder,
		{"duplicates removed", []int64{3, 1, 3, 2, 1placeholder, []int64{1, 2, 3placeholderplaceholder,
		{"zero filtered", []int64{0, 1, 2placeholder, []int64{1, 2placeholderplaceholder,
		{"negative filtered", []int64{-5, -1, 3placeholder, []int64{3placeholderplaceholder,
		{"all invalid", []int64{0, -1, -2placeholder, []int64{placeholderplaceholder,
		{"sorted output", []int64{9, 3, 7, 1placeholder, []int64{1, 3, 7, 9placeholderplaceholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeInt64IDList(tc.in)
			if tc.want == nil {
				require.Nil(t, got)
		placeholder else {
				require.Equal(t, tc.want, got)
		placeholder
	placeholder)
placeholder
placeholder

func TestBuildAccountTodayStatsBatchCacheKey(t *testing.T) {
	tests := []struct {
		name string
		ids  []int64
		want string
placeholder{
		{"empty", nil, "accounts_today_stats_empty"placeholder,
		{"single", []int64{42placeholder, "accounts_today_stats:42"placeholder,
		{"multiple", []int64{1, 2, 3placeholder, "accounts_today_stats:1,2,3"placeholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := buildAccountTodayStatsBatchCacheKey(tc.ids)
			require.Equal(t, tc.want, got)
	placeholder)
placeholder
placeholder
