package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

const liveAttestationHeader = "x-oai-attestation"

type liveAttestationAES struct {
	key [32]byte
placeholder

func newLiveAttestationCipher(cfg *config.Config) SecretEncryptor {
	if cfg == nil || strings.TrimSpace(cfg.JWT.Secret) == "" {
		return nil
placeholder
	return &liveAttestationAES{
		key: sha256.Sum256([]byte("sub2api/live-attestation/v1\x00" + cfg.JWT.Secret)),
placeholder
placeholder

func (c *liveAttestationAES) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(c.key[:])
	if err != nil {
		return "", err
placeholder
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
placeholder
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate Live attestation nonce: %w", err)
placeholder
	encrypted := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.RawStdEncoding.EncodeToString(encrypted), nil
placeholder

func (c *liveAttestationAES) Decrypt(ciphertext string) (string, error) {
	encrypted, err := base64.RawStdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("decode Live attestation: %w", err)
placeholder
	block, err := aes.NewCipher(c.key[:])
	if err != nil {
		return "", err
placeholder
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
placeholder
	if len(encrypted) < gcm.NonceSize() {
		return "", errors.New("encrypted Live attestation is too short")
placeholder
	plaintext, err := gcm.Open(nil, encrypted[:gcm.NonceSize()], encrypted[gcm.NonceSize():], nil)
	if err != nil {
		return "", fmt.Errorf("decrypt Live attestation: %w", err)
placeholder
	return string(plaintext), nil
placeholder

func (s *OpenAIGatewayService) prepareLiveAttestation(ctx context.Context) (string, string, error) {
	if s == nil || s.liveAttestation == nil {
		return "", "", &LiveAttestationUnavailableError{
			Reason: "Sub2API has no platform DeviceCheck provider",
	placeholder
placeholder
	if s.liveAttestationCipher == nil {
		return "", "", &LiveAttestationUnavailableError{
			Reason: "JWT secret is required to protect the Sideband attestation",
	placeholder
placeholder
	header, err := s.liveAttestation.Generate(ctx)
	if err != nil {
		return "", "", &LiveAttestationUnavailableError{Reason: err.Error()placeholder
placeholder
	ciphertext, err := s.liveAttestationCipher.Encrypt(header)
	if err != nil {
		return "", "", &LiveAttestationUnavailableError{
			Reason: "failed to protect the generated DeviceCheck attestation",
	placeholder
placeholder
	return header, ciphertext, nil
placeholder

func (s *OpenAIGatewayService) decryptLiveAttestation(record *LiveCallRecord) (string, error) {
	if record == nil || strings.TrimSpace(record.AttestationCiphertext) == "" || s.liveAttestationCipher == nil {
		return "", &LiveAttestationUnavailableError{
			Reason: "the Live call has no reusable DeviceCheck attestation",
	placeholder
placeholder
	header, err := s.liveAttestationCipher.Decrypt(record.AttestationCiphertext)
	if err != nil {
		return "", &LiveAttestationUnavailableError{
			Reason: "the Live call DeviceCheck attestation cannot be decrypted on this instance",
	placeholder
placeholder
	return header, nil
placeholder
