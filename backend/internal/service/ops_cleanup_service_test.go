package service

import (
	"testing"
	"time"
)

func TestOpsCleanupPlan(t *testing.T) {
	now := time.Date(2026, 4, 29, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name         string
		days         int
		wantOK       bool
		wantTruncate bool
		wantCutoff   time.Time
placeholder{
		{name: "negative skips", days: -1, wantOK: falseplaceholder,
		{name: "zero truncates", days: 0, wantOK: true, wantTruncate: trueplaceholder,
		{name: "positive yields past cutoff", days: 7, wantOK: true, wantCutoff: now.AddDate(0, 0, -7)placeholder,
placeholder

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cutoff, truncate, ok := opsCleanupPlan(now, tc.days)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
		placeholder
			if !ok {
				return
		placeholder
			if truncate != tc.wantTruncate {
				t.Fatalf("truncate = %v, want %v", truncate, tc.wantTruncate)
		placeholder
			if !tc.wantTruncate && !cutoff.Equal(tc.wantCutoff) {
				t.Fatalf("cutoff = %v, want %v", cutoff, tc.wantCutoff)
		placeholder
	placeholder)
placeholder
placeholder

func TestIsMissingRelationError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
placeholder{
		{name: "nil is not missing", err: nil, want: falseplaceholder,
		{name: "match relation does not exist", err: fakeErr(`pq: relation "ops_error_logs" does not exist`), want: trueplaceholder,
		{name: "match case-insensitive", err: fakeErr(`ERROR: Relation "x" Does Not Exist`), want: trueplaceholder,
		{name: "non-matching error", err: fakeErr("connection refused"), want: falseplaceholder,
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isMissingRelationError(tc.err); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
		placeholder
	placeholder)
placeholder
placeholder

type fakeErr string

func (e fakeErr) Error() string { return string(e) placeholder
