//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type geminiTokenCacheStub struct {
	deletedKeys []string
	deleteErr   error
placeholder

func (s *geminiTokenCacheStub) GetAccessToken(ctx context.Context, cacheKey string) (string, error) {
	return "", nil
placeholder

func (s *geminiTokenCacheStub) SetAccessToken(ctx context.Context, cacheKey string, token string, ttl time.Duration) error {
	return nil
placeholder

func (s *geminiTokenCacheStub) DeleteAccessToken(ctx context.Context, cacheKey string) error {
	s.deletedKeys = append(s.deletedKeys, cacheKey)
	return s.deleteErr
placeholder

func (s *geminiTokenCacheStub) AcquireRefreshLock(ctx context.Context, cacheKey string, ttl time.Duration) (bool, error) {
	return true, nil
placeholder

func (s *geminiTokenCacheStub) ReleaseRefreshLock(ctx context.Context, cacheKey string) error {
	return nil
placeholder

func TestCompositeTokenCacheInvalidator_Gemini(t *testing.T) {
	cache := &geminiTokenCacheStub{placeholder
	invalidator := NewCompositeTokenCacheInvalidator(cache)
	account := &Account{
		ID:       10,
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"project_id": "project-x",
	placeholder,
placeholder

	err := invalidator.InvalidateToken(context.Background(), account)
placeholder
	require.Equal(t, []string{"project-x"placeholder, cache.deletedKeys)
placeholder

func TestCompositeTokenCacheInvalidator_Antigravity(t *testing.T) {
	cache := &geminiTokenCacheStub{placeholder
	invalidator := NewCompositeTokenCacheInvalidator(cache)
	account := &Account{
		ID:       99,
		Platform: PlatformAntigravity,
		Type:     AccountTypeOAuth,
placeholder
			"project_id": "ag-project",
	placeholder,
placeholder

	err := invalidator.InvalidateToken(context.Background(), account)
placeholder
	require.Equal(t, []string{"ag:ag-project"placeholder, cache.deletedKeys)
placeholder

func TestCompositeTokenCacheInvalidator_SkipNonOAuth(t *testing.T) {
	cache := &geminiTokenCacheStub{placeholder
	invalidator := NewCompositeTokenCacheInvalidator(cache)
	account := &Account{
		ID:       1,
placeholder
		Type:     AccountTypeAPIKey,
placeholder

	err := invalidator.InvalidateToken(context.Background(), account)
placeholder
	require.Empty(t, cache.deletedKeys)
placeholder

func TestCompositeTokenCacheInvalidator_NilCache(t *testing.T) {
	invalidator := NewCompositeTokenCacheInvalidator(nil)
	account := &Account{
		ID:       2,
placeholder
		Type:     AccountTypeOAuth,
placeholder

	err := invalidator.InvalidateToken(context.Background(), account)
placeholder
placeholder
