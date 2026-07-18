package main

import "testing"

func TestHistoricalIngressRejectReason(t *testing.T) {
	tests := []struct {
		name   string
		item   candidate
		reason string
		match  bool
placeholder{
		{name: "standard invalid key", item: candidate{body: `{"code":"INVALID_API_KEY","message":"Invalid API key"placeholder`placeholder, reason: "invalid_key", match: trueplaceholder,
		{name: "google missing key", item: candidate{body: `{"error":{"code":401,"message":"API key is required","status":"UNAUTHENTICATED"placeholderplaceholder`placeholder, reason: "missing_key", match: trueplaceholder,
		{name: "google group deleted", item: candidate{body: `{"error":{"code":403,"message":"API Key 所属分组已删除","status":"PERMISSION_DENIED"placeholderplaceholder`placeholder, reason: "group_deleted", match: trueplaceholder,
		{name: "ip acl", item: candidate{body: `{"code":"ACCESS_DENIED","message":"Access denied. Your IP is 192.0.2.1"placeholder`placeholder, reason: "ip_acl_denied", match: trueplaceholder,
		{name: "user not found remains", item: candidate{body: `{"code":"USER_NOT_FOUND","message":"User associated with API key not found"placeholder`placeholder, match: falseplaceholder,
		{name: "quota remains", item: candidate{body: `{"code":"API_KEY_QUOTA_EXHAUSTED","message":"quota"placeholder`placeholder, match: falseplaceholder,
		{name: "database failure remains", item: candidate{statusCode: 500, message: "Failed to validate API key", body: `{"code":"INTERNAL_ERROR","message":"Failed to validate API key"placeholder`placeholder, match: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reason, ok := historicalIngressRejectReason(tt.item)
			if ok != tt.match || reason != tt.reason {
				t.Fatalf("got (%q, %v), want (%q, %v)", reason, ok, tt.reason, tt.match)
		placeholder
	placeholder)
placeholder
placeholder
