//go:build unit

package repository

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── 测试辅助 ─────────────────────────────────────────────────────────────────

// aesHexKey 构造一个全填充为 b 的 n 字节密钥并以 hex 编码返回。
func aesHexKey(n int, b byte) string {
	raw := make([]byte, n)
	for i := range raw {
		raw[i] = b
placeholder
	return hex.EncodeToString(raw)
placeholder

// aesTestCfg 用给定 hex 密钥字符串构造最小 Config。
func aesTestCfg(keyHex string) *config.Config {
	return &config.Config{
		Totp: config.TotpConfig{EncryptionKey: keyHexplaceholder,
placeholder
placeholder

// aesEncryptor 创建一个持有合法 32 字节密钥的加密器，测试失败时立即终止。
func aesEncryptor(t *testing.T) *AESEncryptor {
placeholder
	enc, err := NewAESEncryptor(aesTestCfg(aesHexKey(32, 0x42)))
placeholder
	require.NotNil(t, enc)
	return enc.(*AESEncryptor)
placeholder

// ── NewAESEncryptor ──────────────────────────────────────────────────────────

func TestNewAESEncryptor_ValidKey32Bytes(t *testing.T) {
	enc, err := NewAESEncryptor(aesTestCfg(aesHexKey(32, 0x01)))
placeholder
	require.NotNil(t, enc)
placeholder

// 16 / 24 字节密钥在 AES 体系内合法，但本实现仅接受 AES-256（32 字节）。
func TestNewAESEncryptor_WrongKeyLength(t *testing.T) {
	tests := []struct {
		name    string
		keySize int
placeholder{
		{"16_bytes_AES128", 16placeholder,
		{"24_bytes_AES192", 24placeholder,
		{"1_byte", 1placeholder,
		{"31_bytes", 31placeholder,
		{"33_bytes", 33placeholder,
		{"64_bytes", 64placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAESEncryptor(aesTestCfg(aesHexKey(tt.keySize, 0x00)))
		placeholder
			assert.Contains(t, err.Error(), "32 bytes")
	placeholder)
placeholder
placeholder

// "配置缺失"场景：空字符串与非法 hex 编码。
func TestNewAESEncryptor_MissingOrInvalidConfig(t *testing.T) {
	tests := []struct {
		name        string
		keyHex      string
		wantContain string
placeholder{
		{"empty_key", "", "32 bytes"placeholder,
		{"invalid_hex_odd_length", "abcde", "invalid totp encryption key"placeholder,
		{"invalid_hex_chars", "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", "invalid totp encryption key"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAESEncryptor(aesTestCfg(tt.keyHex))
		placeholder
			assert.Contains(t, err.Error(), tt.wantContain)
	placeholder)
placeholder
placeholder

// ── 加解密往返（Roundtrip）───────────────────────────────────────────────────

func TestAESEncryptor_RoundTrip(t *testing.T) {
	enc := aesEncryptor(t)

	tests := []struct {
		name      string
		plaintext string
placeholder{
		{"ascii", "Hello, Sub2API!"placeholder,
		{"chinese_multibyte", "你好，世界！这是多字节 UTF-8 文本。"placeholder,
		{"empty_string", ""placeholder,
		{"long_string_gt_1KB", strings.Repeat("x", 2048)placeholder,
		{"special_chars", "!@#$%^&*()_+-=[]{placeholder|;':\",./<>?"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct, err := enc.Encrypt(tt.plaintext)
		placeholder
			require.NotEmpty(t, ct, "密文不应为空（即便明文为空字符串）")

			got, err := enc.Decrypt(ct)
		placeholder
			assert.Equal(t, tt.plaintext, got)
	placeholder)
placeholder
placeholder

// ── IV/Nonce 随机性 ──────────────────────────────────────────────────────────

func TestAESEncryptor_Encrypt_NonceRandomness(t *testing.T) {
	enc := aesEncryptor(t)
	const iterations = 30
	plaintext := "same plaintext for every iteration"

	seen := make(map[string]struct{placeholder, iterations)
	for i := 0; i < iterations; i++ {
		ct, err := enc.Encrypt(plaintext)
	placeholder
		seen[ct] = struct{placeholder{placeholder
placeholder

	// 30 次加密相同明文，每次因随机 Nonce 应产生不同密文。
	assert.Len(t, seen, iterations,
		"每次加密应因随机 Nonce 产生唯一密文，共 %d 次", iterations)
placeholder

// ── Decrypt 错误路径 ──────────────────────────────────────────────────────────

func TestAESDecrypt_InvalidBase64(t *testing.T) {
	enc := aesEncryptor(t)
	_, err := enc.Decrypt("!!!not-valid-base64!!!")
placeholder
	assert.Contains(t, err.Error(), "decode base64")
placeholder

func TestAESDecrypt_TooShort(t *testing.T) {
	enc := aesEncryptor(t)
	// GCM Nonce 为 12 字节；仅提供 2 字节，必然短于 NonceSize。
	short := base64.StdEncoding.EncodeToString([]byte{0x01, 0x02placeholder)
	_, err := enc.Decrypt(short)
placeholder
	assert.Contains(t, err.Error(), "too short")
placeholder

func TestAESDecrypt_TamperedCiphertext(t *testing.T) {
	enc := aesEncryptor(t)

	ct, err := enc.Encrypt("sensitive payload")
placeholder

	raw, err := base64.StdEncoding.DecodeString(ct)
placeholder

	// Nonce 占前 12 字节；翻转其后第一个字节（密文体）。
	raw[12] ^= 0xFF
	_, err = enc.Decrypt(base64.StdEncoding.EncodeToString(raw))
	require.Error(t, err, "篡改密文体后解密应失败")
placeholder

func TestAESDecrypt_TamperedTag(t *testing.T) {
	enc := aesEncryptor(t)

	ct, err := enc.Encrypt("sensitive payload")
placeholder

	raw, err := base64.StdEncoding.DecodeString(ct)
placeholder

	// GCM 认证标签占最后 16 字节；翻转最后一个字节。
	raw[len(raw)-1] ^= 0xFF
	_, err = enc.Decrypt(base64.StdEncoding.EncodeToString(raw))
	require.Error(t, err, "篡改 GCM 标签后解密应失败")
placeholder

// ── 跨实例（Cross-instance）──────────────────────────────────────────────────

func TestAESEncryptor_CrossInstance_SameKey_CanDecrypt(t *testing.T) {
	keyHex := aesHexKey(32, 0xDE)

	enc1, err := NewAESEncryptor(aesTestCfg(keyHex))
placeholder
	enc2, err := NewAESEncryptor(aesTestCfg(keyHex))
placeholder

	plaintext := "cross-instance roundtrip"
	ct, err := enc1.Encrypt(plaintext)
placeholder

	got, err := enc2.Decrypt(ct)
placeholder
	assert.Equal(t, plaintext, got, "相同密钥构造的两个实例应可互相解密")
placeholder

func TestAESEncryptor_CrossInstance_DifferentKey_CannotDecrypt(t *testing.T) {
	enc1, err := NewAESEncryptor(aesTestCfg(aesHexKey(32, 0xAA)))
placeholder
	enc2, err := NewAESEncryptor(aesTestCfg(aesHexKey(32, 0xBB)))
placeholder

	ct, err := enc1.Encrypt("secret message")
placeholder

	_, err = enc2.Decrypt(ct)
	require.Error(t, err, "不同密钥的实例不应能解密对方的密文")
placeholder
