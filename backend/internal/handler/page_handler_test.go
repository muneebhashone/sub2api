package handler

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanPageImageRelativePath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
		ok   bool
placeholder{
		{name: "single filename", in: "logo.png", want: "logo.png", ok: trueplaceholder,
		{name: "nested path", in: "images/logo.png", want: filepath.Join("images", "logo.png"), ok: trueplaceholder,
		{name: "dot prefix", in: "./logo.png", want: "logo.png", ok: trueplaceholder,
		{name: "url escaped slash", in: "images%2Flogo.png", want: filepath.Join("images", "logo.png"), ok: trueplaceholder,
		{name: "parent traversal", in: "../secret.png", ok: falseplaceholder,
		{name: "encoded parent traversal", in: "%2e%2e/secret.png", ok: falseplaceholder,
		{name: "backslash traversal", in: `images\secret.png`, ok: falseplaceholder,
		{name: "absolute path", in: "/etc/passwd", ok: falseplaceholder,
		{name: "encoded absolute path", in: "%2fetc/passwd", ok: falseplaceholder,
		{name: "encoded nul byte", in: "logo.png%00", ok: falseplaceholder,
		{name: "invalid escape", in: "logo.png%zz", ok: falseplaceholder,
		{name: "empty path", in: "", ok: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := cleanPageImageRelativePath(tt.in)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
		placeholder
			if got != tt.want {
				t.Fatalf("path = %q, want %q", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestResolvePageImagePath(t *testing.T) {
	root := t.TempDir()
	pagesDir := filepath.Join(root, "pages")
	base := filepath.Join(pagesDir, "guide")
	if err := os.MkdirAll(filepath.Join(base, "images"), 0755); err != nil {
		t.Fatalf("create images dir: %v", err)
placeholder
	if err := os.WriteFile(filepath.Join(base, "logo.png"), []byte("fake"), 0644); err != nil {
		t.Fatalf("create direct image: %v", err)
placeholder
	if err := os.WriteFile(filepath.Join(base, "images", "logo.png"), []byte("fake"), 0644); err != nil {
		t.Fatalf("create image: %v", err)
placeholder

	got, ok := resolvePageImagePath(pagesDir, base, "logo.png")
	if !ok {
		t.Fatal("expected direct image path to be accepted")
placeholder
	want := mustEvalSymlinks(t, filepath.Join(base, "logo.png"))
	if got != want {
		t.Fatalf("path = %q, want %q", got, want)
placeholder

	got, ok = resolvePageImagePath(pagesDir, base, "images/logo.png")
	if !ok {
		t.Fatal("expected nested image path to be accepted")
placeholder
	want = mustEvalSymlinks(t, filepath.Join(base, "images", "logo.png"))
	if got != want {
		t.Fatalf("path = %q, want %q", got, want)
placeholder

	if got, ok := resolvePageImagePath(pagesDir, base, "../guide.md"); ok {
		t.Fatalf("expected traversal to be rejected, got %q", got)
placeholder
placeholder

func TestResolvePageImagePathRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	pagesDir := filepath.Join(root, "pages")
	base := filepath.Join(pagesDir, "guide")
	outside := filepath.Join(root, "outside")

	if err := os.MkdirAll(base, 0755); err != nil {
		t.Fatalf("create page dir: %v", err)
placeholder
	if err := os.MkdirAll(outside, 0755); err != nil {
		t.Fatalf("create outside dir: %v", err)
placeholder
	if err := os.WriteFile(filepath.Join(outside, "secret.png"), []byte("secret"), 0644); err != nil {
		t.Fatalf("create outside file: %v", err)
placeholder
	if err := os.Symlink(outside, filepath.Join(base, "images")); err != nil {
		t.Skipf("symlink not supported: %v", err)
placeholder

	if got, ok := resolvePageImagePath(pagesDir, base, "images/secret.png"); ok {
		t.Fatalf("expected symlink escape to be rejected, got %q", got)
placeholder
placeholder

func mustEvalSymlinks(t *testing.T, path string) string {
placeholder

	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("eval symlinks for %q: %v", path, err)
placeholder
	return realPath
placeholder
