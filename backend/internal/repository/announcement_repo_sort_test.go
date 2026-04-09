package repository

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

func TestAnnouncementListOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params pagination.PaginationParams
		wantBy string
		want   string
placeholder{
		{
			name:   "default created_at desc",
			params: pagination.PaginationParams{placeholder,
			wantBy: "created_at",
			want:   "desc",
	placeholder,
		{
			name: "title asc",
			params: pagination.PaginationParams{
				SortBy:    "title",
				SortOrder: "ASC",
		placeholder,
			wantBy: "title",
			want:   "asc",
	placeholder,
		{
			name: "status desc",
			params: pagination.PaginationParams{
				SortBy:    "status",
				SortOrder: "desc",
		placeholder,
			wantBy: "status",
			want:   "desc",
	placeholder,
		{
			name: "invalid falls back",
			params: pagination.PaginationParams{
				SortBy:    "sideways",
				SortOrder: "wat",
		placeholder,
			wantBy: "created_at",
			want:   "desc",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotBy, gotOrder := announcementListOrder(tt.params)
			if gotBy != tt.wantBy || gotOrder != tt.want {
				t.Fatalf("announcementListOrder(%+v) = (%q, %q), want (%q, %q)", tt.params, gotBy, gotOrder, tt.wantBy, tt.want)
		placeholder
	placeholder)
placeholder
placeholder
