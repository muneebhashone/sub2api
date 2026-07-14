package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortedUniqueAccountIDs(t *testing.T) {
	tests := []struct {
		name  string
		input []int64
		want  []int64
placeholder{
		{name: "unsorted duplicates", input: []int64{12, 3, 12, 8, 3placeholder, want: []int64{3, 8, 12placeholderplaceholder,
		{name: "already sorted", input: []int64{3, 8, 12placeholder, want: []int64{3, 8, 12placeholderplaceholder,
		{name: "single", input: []int64{3placeholder, want: []int64{3placeholderplaceholder,
		{name: "empty", input: []int64{placeholder, want: []int64{placeholderplaceholder,
		{name: "nil", input: nil, want: nilplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, sortedUniqueAccountIDs(tt.input))
	placeholder)
placeholder
placeholder
