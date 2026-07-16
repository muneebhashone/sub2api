package middleware

import "testing"

func TestDeriveAuditAction(t *testing.T) {
	cases := []struct {
		method string
		path   string
		want   string
placeholder{
		{"PUT", "/api/v1/admin/accounts/:id", "admin.accounts.update"placeholder,
		{"POST", "/api/v1/admin/accounts", "admin.accounts.create"placeholder,
		{"DELETE", "/api/v1/admin/backups/:id", "admin.backups.delete"placeholder,
		{"GET", "/api/v1/admin/users/:id/api-keys", "admin.users.api_keys.read"placeholder,
		{"POST", "/api/v1/admin/redeem-codes/batch", "admin.redeem_codes.batch.create"placeholder,
placeholder
	for _, tc := range cases {
		if got := deriveAuditAction(tc.method, tc.path); got != tc.want {
			t.Fatalf("deriveAuditAction(%q, %q) = %q, want %q", tc.method, tc.path, got, tc.want)
	placeholder
placeholder
placeholder
