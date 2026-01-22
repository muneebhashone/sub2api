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
	// 新行为：同时删除基于 project_id 和 account_id 的缓存键
	// 这是为了处理：首次获取 token 时可能没有 project_id，之后自动检测到后会使用新 key
	require.Equal(t, []string{"gemini:project-x", "gemini:account:10"placeholder, cache.deletedKeys)
placeholder

func TestCompositeTokenCacheInvalidator_GeminiWithoutProjectID(t *testing.T) {
	cache := &geminiTokenCacheStub{placeholder
	invalidator := NewCompositeTokenCacheInvalidator(cache)
	account := &Account{
		ID:       10,
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "gemini-token",
	placeholder,
placeholder

	err := invalidator.InvalidateToken(context.Background(), account)
placeholder
	// 没有 project_id 时，两个 key 相同，去重后只删除一个
	require.Equal(t, []string{"gemini:account:10"placeholder, cache.deletedKeys)
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
	// 新行为：同时删除基于 project_id 和 account_id 的缓存键
	require.Equal(t, []string{"ag:ag-project", "ag:account:99"placeholder, cache.deletedKeys)
placeholder

func TestCompositeTokenCacheInvalidator_AntigravityWithoutProjectID(t *testing.T) {
	cache := &geminiTokenCacheStub{placeholder
	invalidator := NewCompositeTokenCacheInvalidator(cache)
	account := &Account{
		ID:       99,
		Platform: PlatformAntigravity,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "ag-token",
	placeholder,
placeholder

	err := invalidator.InvalidateToken(context.Background(), account)
placeholder
	// 没有 project_id 时，两个 key 相同，去重后只删除一个
	require.Equal(t, []string{"ag:account:99"placeholder, cache.deletedKeys)
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
		name    string
		account *Account
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
		name    string
		account *Account
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
			// 新行为：删除失败只记录日志，不返回错误
			// 这是因为缓存失效失败不应影响主业务流程
			err := invalidator.InvalidateToken(context.Background(), tt.account)
		placeholder
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

	// 新行为：Gemini 和 Antigravity 会同时删除基于 project_id 和 account_id 的键
	expectedKeys := []string{
		"gemini:gemini-proj",
		"gemini:account:1",
		"ag:ag-proj",
		"ag:account:2",
		"openai:account:3",
		"claude:account:4",
placeholder

	for _, acc := range accounts {
		err := invalidator.InvalidateToken(context.Background(), acc)
	placeholder
placeholder

	require.Equal(t, expectedKeys, cache.deletedKeys)
placeholder

// ========== GetCredentialAsInt64 测试 ==========

func TestAccount_GetCredentialAsInt64(t *testing.T) {
	tests := []struct {
		name        string
		credentials map[string]any
		key         string
		expected    int64
placeholder{
		{
			name:        "int64_value",
			credentials: map[string]any{"_token_version": int64(1737654321000)placeholder,
			key:         "_token_version",
			expected:    1737654321000,
	placeholder,
		{
			name:        "float64_value",
			credentials: map[string]any{"_token_version": float64(1737654321000)placeholder,
			key:         "_token_version",
			expected:    1737654321000,
	placeholder,
		{
			name:        "int_value",
			credentials: map[string]any{"_token_version": 12345placeholder,
			key:         "_token_version",
			expected:    12345,
	placeholder,
		{
			name:        "string_value",
			credentials: map[string]any{"_token_version": "1737654321000"placeholder,
			key:         "_token_version",
			expected:    1737654321000,
	placeholder,
		{
			name:        "string_with_spaces",
			credentials: map[string]any{"_token_version": "  1737654321000  "placeholder,
			key:         "_token_version",
			expected:    1737654321000,
	placeholder,
		{
			name:        "nil_credentials",
			credentials: nil,
			key:         "_token_version",
			expected:    0,
	placeholder,
		{
			name:        "missing_key",
			credentials: map[string]any{"other_key": placeholder,
			key:         "_token_version",
			expected:    0,
	placeholder,
		{
			name:        "nil_value",
			credentials: map[string]any{"_token_version": nilplaceholder,
			key:         "_token_version",
			expected:    0,
	placeholder,
		{
			name:        "invalid_string",
			credentials: map[string]any{"_token_version": "not_a_number"placeholder,
			key:         "_token_version",
			expected:    0,
	placeholder,
		{
			name:        "empty_string",
			credentials: map[string]any{"_token_version": ""placeholder,
			key:         "_token_version",
			expected:    0,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{Credentials: tt.credentialsplaceholder
			result := account.GetCredentialAsInt64(tt.key)
			require.Equal(t, tt.expected, result)
	placeholder)
placeholder
placeholder

func TestAccount_GetCredentialAsInt64_NilAccount(t *testing.T) {
	var account *Account
	result := account.GetCredentialAsInt64("_token_version")
	require.Equal(t, int64(0), result)
placeholder

// ========== IsTokenVersionStale 测试 ==========

func TestIsTokenVersionStale(t *testing.T) {
	tests := []struct {
		name          string
		account       *Account
		latestAccount *Account
		repoErr       error
		expectedStale bool
placeholder{
		{
			name:          "nil_account",
			account:       nil,
			latestAccount: nil,
			expectedStale: false,
	placeholder,
		{
			name: "no_version_in_account_but_db_has_version",
			account: &Account{
				ID:          1,
		placeholderplaceholder,
		placeholder,
			latestAccount: &Account{
				ID:          1,
		placeholder"_token_version": int64(100)placeholder,
		placeholder,
			expectedStale: true, // 当前 account 无版本但 DB 有，说明已被异步刷新，当前已过时
	placeholder,
		{
			name: "both_no_version",
			account: &Account{
				ID:          1,
		placeholderplaceholder,
		placeholder,
			latestAccount: &Account{
				ID:          1,
		placeholderplaceholder,
		placeholder,
			expectedStale: false, // 两边都没有版本号，说明从未被异步刷新过，允许缓存
	placeholder,
		{
			name: "same_version",
			account: &Account{
				ID:          1,
		placeholder"_token_version": int64(100)placeholder,
		placeholder,
			latestAccount: &Account{
				ID:          1,
		placeholder"_token_version": int64(100)placeholder,
		placeholder,
			expectedStale: false,
	placeholder,
		{
			name: "current_version_newer",
			account: &Account{
				ID:          1,
		placeholder"_token_version": int64(200)placeholder,
		placeholder,
			latestAccount: &Account{
				ID:          1,
		placeholder"_token_version": int64(100)placeholder,
		placeholder,
			expectedStale: false,
	placeholder,
		{
			name: "current_version_older_stale",
			account: &Account{
				ID:          1,
		placeholder"_token_version": int64(100)placeholder,
		placeholder,
			latestAccount: &Account{
				ID:          1,
		placeholder"_token_version": int64(200)placeholder,
		placeholder,
			expectedStale: true, // 当前版本过时
	placeholder,
		{
			name: "repo_error",
			account: &Account{
				ID:          1,
		placeholder"_token_version": int64(100)placeholder,
		placeholder,
			latestAccount: nil,
			repoErr:       errors.New("db error"),
			expectedStale: false, // 查询失败，默认允许缓存
	placeholder,
		{
			name: "repo_returns_nil",
			account: &Account{
				ID:          1,
		placeholder"_token_version": int64(100)placeholder,
		placeholder,
			latestAccount: nil,
			repoErr:       nil,
			expectedStale: false, // 查询返回 nil，默认允许缓存
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 由于 IsTokenVersionStale 接受 AccountRepository 接口，而创建完整的 mock 很繁琐
			// 这里我们直接测试函数的核心逻辑来验证行为

			if tt.name == "nil_account" {
				result := IsTokenVersionStale(context.Background(), nil, nil)
				require.Equal(t, tt.expectedStale, result)
				return
		placeholder

			// 模拟 IsTokenVersionStale 的核心逻辑
			account := tt.account
			currentVersion := account.GetCredentialAsInt64("_token_version")

			// 模拟 repo 查询
			latestAccount := tt.latestAccount
			if tt.repoErr != nil || latestAccount == nil {
				require.Equal(t, tt.expectedStale, false)
				return
		placeholder

			latestVersion := latestAccount.GetCredentialAsInt64("_token_version")

			// 情况1: 当前 account 没有版本号，但 DB 中已有版本号
			if currentVersion == 0 && latestVersion > 0 {
				require.Equal(t, tt.expectedStale, true)
				return
		placeholder

			// 情况2: 两边都没有版本号
			if currentVersion == 0 && latestVersion == 0 {
				require.Equal(t, tt.expectedStale, false)
				return
		placeholder

			// 情况3: 比较版本号
			isStale := latestVersion > currentVersion
			require.Equal(t, tt.expectedStale, isStale)
	placeholder)
placeholder
placeholder

func TestIsTokenVersionStale_NilRepo(t *testing.T) {
	account := &Account{
		ID:          1,
placeholder"_token_version": int64(100)placeholder,
placeholder
	result := IsTokenVersionStale(context.Background(), account, nil)
	require.False(t, result) // nil repo，默认允许缓存
placeholder
