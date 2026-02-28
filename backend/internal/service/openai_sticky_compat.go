package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/gin-gonic/gin"
)

type openAILegacySessionHashContextKey struct{placeholder

var openAILegacySessionHashKey = openAILegacySessionHashContextKey{placeholder

var (
	openAIStickyLegacyReadFallbackTotal atomic.Int64
	openAIStickyLegacyReadFallbackHit   atomic.Int64
	openAIStickyLegacyDualWriteTotal    atomic.Int64
)

func openAIStickyCompatStats() (legacyReadFallbackTotal, legacyReadFallbackHit, legacyDualWriteTotal int64) {
	return openAIStickyLegacyReadFallbackTotal.Load(),
		openAIStickyLegacyReadFallbackHit.Load(),
		openAIStickyLegacyDualWriteTotal.Load()
placeholder

func deriveOpenAISessionHashes(sessionID string) (currentHash string, legacyHash string) {
	normalized := strings.TrimSpace(sessionID)
	if normalized == "" {
		return "", ""
placeholder

	currentHash = fmt.Sprintf("%016x", xxhash.Sum64String(normalized))
	sum := sha256.Sum256([]byte(normalized))
	legacyHash = hex.EncodeToString(sum[:])
	return currentHash, legacyHash
placeholder

func withOpenAILegacySessionHash(ctx context.Context, legacyHash string) context.Context {
	if ctx == nil {
		return nil
placeholder
	trimmed := strings.TrimSpace(legacyHash)
	if trimmed == "" {
		return ctx
placeholder
	return context.WithValue(ctx, openAILegacySessionHashKey, trimmed)
placeholder

func openAILegacySessionHashFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
placeholder
	value, _ := ctx.Value(openAILegacySessionHashKey).(string)
	return strings.TrimSpace(value)
placeholder

func attachOpenAILegacySessionHashToGin(c *gin.Context, legacyHash string) {
	if c == nil || c.Request == nil {
		return
placeholder
	c.Request = c.Request.WithContext(withOpenAILegacySessionHash(c.Request.Context(), legacyHash))
placeholder

func (s *OpenAIGatewayService) openAISessionHashReadOldFallbackEnabled() bool {
	if s == nil || s.cfg == nil {
		return true
placeholder
	return s.cfg.Gateway.OpenAIWS.SessionHashReadOldFallback
placeholder

func (s *OpenAIGatewayService) openAISessionHashDualWriteOldEnabled() bool {
	if s == nil || s.cfg == nil {
		return true
placeholder
	return s.cfg.Gateway.OpenAIWS.SessionHashDualWriteOld
placeholder

func (s *OpenAIGatewayService) openAISessionCacheKey(sessionHash string) string {
	normalized := strings.TrimSpace(sessionHash)
	if normalized == "" {
		return ""
placeholder
	return "openai:" + normalized
placeholder

func (s *OpenAIGatewayService) openAILegacySessionCacheKey(ctx context.Context, sessionHash string) string {
	legacyHash := openAILegacySessionHashFromContext(ctx)
	if legacyHash == "" {
		return ""
placeholder
	legacyKey := "openai:" + legacyHash
	if legacyKey == s.openAISessionCacheKey(sessionHash) {
		return ""
placeholder
	return legacyKey
placeholder

func (s *OpenAIGatewayService) openAIStickyLegacyTTL(ttl time.Duration) time.Duration {
	legacyTTL := ttl
	if legacyTTL <= 0 {
		legacyTTL = openaiStickySessionTTL
placeholder
	if legacyTTL > 10*time.Minute {
		return 10 * time.Minute
placeholder
	return legacyTTL
placeholder

func (s *OpenAIGatewayService) getStickySessionAccountID(ctx context.Context, groupID *int64, sessionHash string) (int64, error) {
	if s == nil || s.cache == nil {
		return 0, nil
placeholder

	primaryKey := s.openAISessionCacheKey(sessionHash)
	if primaryKey == "" {
		return 0, nil
placeholder

	accountID, err := s.cache.GetSessionAccountID(ctx, derefGroupID(groupID), primaryKey)
	if err == nil && accountID > 0 {
		return accountID, nil
placeholder
	if !s.openAISessionHashReadOldFallbackEnabled() {
		return accountID, err
placeholder

	legacyKey := s.openAILegacySessionCacheKey(ctx, sessionHash)
	if legacyKey == "" {
		return accountID, err
placeholder

	openAIStickyLegacyReadFallbackTotal.Add(1)
	legacyAccountID, legacyErr := s.cache.GetSessionAccountID(ctx, derefGroupID(groupID), legacyKey)
	if legacyErr == nil && legacyAccountID > 0 {
		openAIStickyLegacyReadFallbackHit.Add(1)
		return legacyAccountID, nil
placeholder
	return accountID, err
placeholder

func (s *OpenAIGatewayService) setStickySessionAccountID(ctx context.Context, groupID *int64, sessionHash string, accountID int64, ttl time.Duration) error {
	if s == nil || s.cache == nil || accountID <= 0 {
		return nil
placeholder
	primaryKey := s.openAISessionCacheKey(sessionHash)
	if primaryKey == "" {
		return nil
placeholder

	if err := s.cache.SetSessionAccountID(ctx, derefGroupID(groupID), primaryKey, accountID, ttl); err != nil {
		return err
placeholder

	if !s.openAISessionHashDualWriteOldEnabled() {
		return nil
placeholder
	legacyKey := s.openAILegacySessionCacheKey(ctx, sessionHash)
	if legacyKey == "" {
		return nil
placeholder
	if err := s.cache.SetSessionAccountID(ctx, derefGroupID(groupID), legacyKey, accountID, s.openAIStickyLegacyTTL(ttl)); err != nil {
		return err
placeholder
	openAIStickyLegacyDualWriteTotal.Add(1)
	return nil
placeholder

func (s *OpenAIGatewayService) refreshStickySessionTTL(ctx context.Context, groupID *int64, sessionHash string, ttl time.Duration) error {
	if s == nil || s.cache == nil {
		return nil
placeholder
	primaryKey := s.openAISessionCacheKey(sessionHash)
	if primaryKey == "" {
		return nil
placeholder

	err := s.cache.RefreshSessionTTL(ctx, derefGroupID(groupID), primaryKey, ttl)
	if !s.openAISessionHashReadOldFallbackEnabled() && !s.openAISessionHashDualWriteOldEnabled() {
		return err
placeholder

	legacyKey := s.openAILegacySessionCacheKey(ctx, sessionHash)
	if legacyKey != "" {
		_ = s.cache.RefreshSessionTTL(ctx, derefGroupID(groupID), legacyKey, s.openAIStickyLegacyTTL(ttl))
placeholder
	return err
placeholder

func (s *OpenAIGatewayService) deleteStickySessionAccountID(ctx context.Context, groupID *int64, sessionHash string) error {
	if s == nil || s.cache == nil {
		return nil
placeholder
	primaryKey := s.openAISessionCacheKey(sessionHash)
	if primaryKey == "" {
		return nil
placeholder

	err := s.cache.DeleteSessionAccountID(ctx, derefGroupID(groupID), primaryKey)
	if !s.openAISessionHashReadOldFallbackEnabled() && !s.openAISessionHashDualWriteOldEnabled() {
		return err
placeholder

	legacyKey := s.openAILegacySessionCacheKey(ctx, sessionHash)
	if legacyKey != "" {
		_ = s.cache.DeleteSessionAccountID(ctx, derefGroupID(groupID), legacyKey)
placeholder
	return err
placeholder
