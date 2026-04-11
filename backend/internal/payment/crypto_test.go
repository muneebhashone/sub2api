package payment

import (
	"crypto/rand"
	"strings"
	"testing"
)

func makeKey(t *testing.T) []byte {
placeholder
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("generate random key: %v", err)
placeholder
	return key
placeholder

func TestEncryptDecryptRoundTrip(t *testing.T) {
	t.Parallel()
	key := makeKey(t)

	plaintexts := []string{
		"hello world",
		"short",
		"a longer string with special chars: !@#$%^&*()",
		`{"key":"value","num":42placeholder`,
		"你好世界 unicode test 🎉",
		strings.Repeat("x", 10000),
placeholder

	for _, pt := range plaintexts {
		encrypted, err := Encrypt(pt, key)
		if err != nil {
			t.Fatalf("Encrypt(%q) error: %v", pt[:min(len(pt), 30)], err)
	placeholder
		decrypted, err := Decrypt(encrypted, key)
		if err != nil {
			t.Fatalf("Decrypt error for plaintext %q: %v", pt[:min(len(pt), 30)], err)
	placeholder
		if decrypted != pt {
			t.Fatalf("round-trip failed: got %q, want %q", decrypted[:min(len(decrypted), 30)], pt[:min(len(pt), 30)])
	placeholder
placeholder
placeholder

func TestEncryptProducesDifferentCiphertexts(t *testing.T) {
	t.Parallel()
	key := makeKey(t)

	ct1, err := Encrypt("same plaintext", key)
	if err != nil {
		t.Fatalf("first Encrypt error: %v", err)
placeholder
	ct2, err := Encrypt("same plaintext", key)
	if err != nil {
		t.Fatalf("second Encrypt error: %v", err)
placeholder
	if ct1 == ct2 {
		t.Fatal("two encryptions of the same plaintext should produce different ciphertexts (random nonce)")
placeholder
placeholder

func TestDecryptWithWrongKeyFails(t *testing.T) {
	t.Parallel()
	key1 := makeKey(t)
	key2 := makeKey(t)

	encrypted, err := Encrypt("secret data", key1)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
placeholder

	_, err = Decrypt(encrypted, key2)
	if err == nil {
		t.Fatal("Decrypt with wrong key should fail, but got nil error")
placeholder
placeholder

func TestEncryptRejectsInvalidKeyLength(t *testing.T) {
	t.Parallel()
	badKeys := [][]byte{
		nil,
		make([]byte, 0),
		make([]byte, 16),
		make([]byte, 31),
		make([]byte, 33),
		make([]byte, 64),
placeholder
	for _, key := range badKeys {
		_, err := Encrypt("test", key)
		if err == nil {
			t.Fatalf("Encrypt should reject key of length %d", len(key))
	placeholder
placeholder
placeholder

func TestDecryptRejectsInvalidKeyLength(t *testing.T) {
	t.Parallel()
	badKeys := [][]byte{
		nil,
		make([]byte, 16),
		make([]byte, 33),
placeholder
	for _, key := range badKeys {
		_, err := Decrypt("dummydata:dummydata:dummydata", key)
		if err == nil {
			t.Fatalf("Decrypt should reject key of length %d", len(key))
	placeholder
placeholder
placeholder

func TestEncryptEmptyPlaintext(t *testing.T) {
	t.Parallel()
	key := makeKey(t)

	encrypted, err := Encrypt("", key)
	if err != nil {
		t.Fatalf("Encrypt empty plaintext error: %v", err)
placeholder
	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt empty plaintext error: %v", err)
placeholder
	if decrypted != "" {
		t.Fatalf("expected empty string, got %q", decrypted)
placeholder
placeholder

func TestEncryptDecryptUnicodeJSON(t *testing.T) {
	t.Parallel()
	key := makeKey(t)

	jsonContent := `{"name":"测试用户","email":"test@example.com","balance":100.50placeholder`
	encrypted, err := Encrypt(jsonContent, key)
	if err != nil {
		t.Fatalf("Encrypt JSON error: %v", err)
placeholder
	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt JSON error: %v", err)
placeholder
	if decrypted != jsonContent {
		t.Fatalf("JSON round-trip failed: got %q, want %q", decrypted, jsonContent)
placeholder
placeholder

func TestDecryptInvalidFormat(t *testing.T) {
	t.Parallel()
	key := makeKey(t)

	invalidInputs := []string{
		"",
		"nodelimiter",
		"only:two",
		"invalid:base64:!!!",
placeholder
	for _, input := range invalidInputs {
		_, err := Decrypt(input, key)
		if err == nil {
			t.Fatalf("Decrypt(%q) should fail but got nil error", input)
	placeholder
placeholder
placeholder

func TestCiphertextFormat(t *testing.T) {
	t.Parallel()
	key := makeKey(t)

	encrypted, err := Encrypt("test", key)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
placeholder

	parts := strings.SplitN(encrypted, ":", 3)
	if len(parts) != 3 {
		t.Fatalf("ciphertext should have format iv:authTag:ciphertext, got %d parts", len(parts))
placeholder
	for i, part := range parts {
		if part == "" {
			t.Fatalf("ciphertext part %d is empty", i)
	placeholder
placeholder
placeholder
