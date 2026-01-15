//go:build unit

package service

import (
	"context"
	"errors"
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

func TestCompositeTokenCacheInvalidator_OpenAI(t *testing.T) {
	cache := &geminiTokenCacheStub{placeholder
	invalidator := NewCompositeTokenCacheInvalidator(cache)
	account := &Account{
		ID:       500,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "openai-token",
	placeholder,
placeholder

	err := invalidator.InvalidateToken(context.Background(), account)
placeholder
	require.Equal(t, []string{"openai:account:500"placeholder, cache.deletedKeys)
placeholder

func TestCompositeTokenCacheInvalidator_Claude(t *testing.T) {
	cache := &geminiTokenCacheStub{placeholder
	invalidator := NewCompositeTokenCacheInvalidator(cache)
	account := &Account{
		ID:       600,
		Platform: PlatformAnthropic,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "claude-token",
	placeholder,
placeholder

	err := invalidator.InvalidateToken(context.Background(), account)
placeholder
	require.Equal(t, []string{"claude:account:600"placeholder, cache.deletedKeys)
placeholder

func TestCompositeTokenCacheInvalidator_SkipNonOAuth(t *testing.T) {
	cache := &geminiTokenCacheStub{placeholder
	invalidator := NewCompositeTokenCacheInvalidator(cache)

	tests := []struct {
		name     string
		account  *Account
placeholder{
		{
			name: "gemini_api_key",
			account: &Account{
				ID:       1,
		placeholder
				Type:     AccountTypeAPIKey,
		placeholder,
	placeholder,
		{
			name: "openai_api_key",
			account: &Account{
				ID:       2,
				Platform: PlatformOpenAI,
				Type:     AccountTypeAPIKey,
		placeholder,
	placeholder,
		{
			name: "claude_api_key",
			account: &Account{
				ID:       3,
				Platform: PlatformAnthropic,
				Type:     AccountTypeAPIKey,
		placeholder,
	placeholder,
		{
			name: "claude_setup_token",
			account: &Account{
				ID:       4,
				Platform: PlatformAnthropic,
				Type:     AccountTypeSetupToken,
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.deletedKeys = nil
			err := invalidator.InvalidateToken(context.Background(), tt.account)
		placeholder
			require.Empty(t, cache.deletedKeys)
	placeholder)
placeholder
placeholder

func TestCompositeTokenCacheInvalidator_SkipUnsupportedPlatform(t *testing.T) {
	cache := &geminiTokenCacheStub{placeholder
	invalidator := NewCompositeTokenCacheInvalidator(cache)
	account := &Account{
		ID:       100,
		Platform: "unknown-platform",
		Type:     AccountTypeOAuth,
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

func TestCompositeTokenCacheInvalidator_NilAccount(t *testing.T) {
	cache := &geminiTokenCacheStub{placeholder
	invalidator := NewCompositeTokenCacheInvalidator(cache)

	err := invalidator.InvalidateToken(context.Background(), nil)
placeholder
	require.Empty(t, cache.deletedKeys)
placeholder

func TestCompositeTokenCacheInvalidator_NilInvalidator(t *testing.T) {
	var invalidator *CompositeTokenCacheInvalidator
	account := &Account{
		ID:       5,
placeholder
		Type:     AccountTypeOAuth,
placeholder

	err := invalidator.InvalidateToken(context.Background(), account)
placeholder
placeholder

func TestCompositeTokenCacheInvalidator_DeleteError(t *testing.T) {
	expectedErr := errors.New("redis connection failed")
	cache := &geminiTokenCacheStub{deleteErr: expectedErrplaceholder
	invalidator := NewCompositeTokenCacheInvalidator(cache)

	tests := []struct {
		name     string
		account  *Account
placeholder{
		{
			name: "openai_delete_error",
			account: &Account{
				ID:       700,
				Platform: PlatformOpenAI,
				Type:     AccountTypeOAuth,
		placeholder,
	placeholder,
		{
			name: "claude_delete_error",
			account: &Account{
				ID:       800,
				Platform: PlatformAnthropic,
				Type:     AccountTypeOAuth,
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := invalidator.InvalidateToken(context.Background(), tt.account)
		placeholder
			require.Equal(t, expectedErr, err)
	placeholder)
placeholder
placeholder

func TestCompositeTokenCacheInvalidator_AllPlatformsIntegration(t *testing.T) {
	// 测试所有平台的缓存键生成和删除
	cache := &geminiTokenCacheStub{placeholder
	invalidator := NewCompositeTokenCacheInvalidator(cache)

	accounts := []*Account{
		{ID: 1, Platform: PlatformGemini, Type: AccountTypeOAuth, Credentials: map[string]any{"project_id": "gemini-proj"placeholderplaceholder,
		{ID: 2, Platform: PlatformAntigravity, Type: AccountTypeOAuth, Credentials: map[string]any{"project_id": "ag-proj"placeholderplaceholder,
		{ID: 3, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder,
		{ID: 4, Platform: PlatformAnthropic, Type: AccountTypeOAuthplaceholder,
placeholder

	expectedKeys := []string{
		"gemini-proj",
		"ag:ag-proj",
		"openai:account:3",
		"claude:account:4",
placeholder

	for _, acc := range accounts {
		err := invalidator.InvalidateToken(context.Background(), acc)
	placeholder
placeholder

	require.Equal(t, expectedKeys, cache.deletedKeys)
placeholder
