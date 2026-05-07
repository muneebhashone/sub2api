package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type contentModerationTestSettingRepo struct {
	values map[string]string
placeholder

func (r *contentModerationTestSettingRepo) Get(ctx context.Context, key string) (*Setting, error) {
	if value, ok := r.values[key]; ok {
		return &Setting{Key: key, Value: valueplaceholder, nil
placeholder
	return nil, ErrSettingNotFound
placeholder

func (r *contentModerationTestSettingRepo) GetValue(ctx context.Context, key string) (string, error) {
	if value, ok := r.values[key]; ok {
		return value, nil
placeholder
	return "", ErrSettingNotFound
placeholder

func (r *contentModerationTestSettingRepo) Set(ctx context.Context, key, value string) error {
	if r.values == nil {
		r.values = map[string]string{placeholder
placeholder
	r.values[key] = value
	return nil
placeholder

func (r *contentModerationTestSettingRepo) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	out := map[string]string{placeholder
	for _, key := range keys {
		if value, ok := r.values[key]; ok {
			out[key] = value
	placeholder
placeholder
	return out, nil
placeholder

func (r *contentModerationTestSettingRepo) SetMultiple(ctx context.Context, settings map[string]string) error {
	if r.values == nil {
		r.values = map[string]string{placeholder
placeholder
	for key, value := range settings {
		r.values[key] = value
placeholder
	return nil
placeholder

func (r *contentModerationTestSettingRepo) GetAll(ctx context.Context) (map[string]string, error) {
	out := make(map[string]string, len(r.values))
	for key, value := range r.values {
		out[key] = value
placeholder
	return out, nil
placeholder

func (r *contentModerationTestSettingRepo) Delete(ctx context.Context, key string) error {
	delete(r.values, key)
	return nil
placeholder

type contentModerationTestRepo struct {
	logs []ContentModerationLog
placeholder

func (r *contentModerationTestRepo) CreateLog(ctx context.Context, log *ContentModerationLog) error {
	if log != nil {
		r.logs = append(r.logs, *log)
placeholder
	return nil
placeholder

func (r *contentModerationTestRepo) ListLogs(ctx context.Context, filter ContentModerationLogFilter) ([]ContentModerationLog, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder

func (r *contentModerationTestRepo) CountFlaggedByUserSince(ctx context.Context, userID int64, since time.Time) (int, error) {
	return 0, nil
placeholder

func (r *contentModerationTestRepo) CleanupExpiredLogs(ctx context.Context, hitBefore time.Time, nonHitBefore time.Time) (*ContentModerationCleanupResult, error) {
	return &ContentModerationCleanupResult{placeholder, nil
placeholder

type contentModerationTestHashCache struct {
	hashes        map[string]struct{placeholder
	recorded      []string
	checked       []string
	deleted       []string
	hasResult     bool
	hasResultUsed bool
placeholder

type contentModerationTestUserRepo struct {
	user    *User
	updated []User
placeholder

func (r *contentModerationTestUserRepo) Create(ctx context.Context, user *User) error {
	panic("unexpected Create call")
placeholder

func (r *contentModerationTestUserRepo) GetByID(ctx context.Context, id int64) (*User, error) {
	if r.user == nil {
		return nil, ErrUserNotFound
placeholder
	clone := *r.user
	return &clone, nil
placeholder

func (r *contentModerationTestUserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	panic("unexpected GetByEmail call")
placeholder

func (r *contentModerationTestUserRepo) GetFirstAdmin(ctx context.Context) (*User, error) {
	panic("unexpected GetFirstAdmin call")
placeholder

func (r *contentModerationTestUserRepo) Update(ctx context.Context, user *User) error {
	if user == nil {
		return nil
placeholder
	clone := *user
	r.updated = append(r.updated, clone)
	r.user = &clone
	return nil
placeholder

func (r *contentModerationTestUserRepo) Delete(ctx context.Context, id int64) error {
	panic("unexpected Delete call")
placeholder

func (r *contentModerationTestUserRepo) GetUserAvatar(ctx context.Context, userID int64) (*UserAvatar, error) {
	panic("unexpected GetUserAvatar call")
placeholder

func (r *contentModerationTestUserRepo) UpsertUserAvatar(ctx context.Context, userID int64, input UpsertUserAvatarInput) (*UserAvatar, error) {
	panic("unexpected UpsertUserAvatar call")
placeholder

func (r *contentModerationTestUserRepo) DeleteUserAvatar(ctx context.Context, userID int64) error {
	panic("unexpected DeleteUserAvatar call")
placeholder

func (r *contentModerationTestUserRepo) List(ctx context.Context, params pagination.PaginationParams) ([]User, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (r *contentModerationTestUserRepo) ListWithFilters(ctx context.Context, params pagination.PaginationParams, filters UserListFilters) ([]User, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
placeholder

func (r *contentModerationTestUserRepo) GetLatestUsedAtByUserIDs(ctx context.Context, userIDs []int64) (map[int64]*time.Time, error) {
	panic("unexpected GetLatestUsedAtByUserIDs call")
placeholder

func (r *contentModerationTestUserRepo) GetLatestUsedAtByUserID(ctx context.Context, userID int64) (*time.Time, error) {
	panic("unexpected GetLatestUsedAtByUserID call")
placeholder

func (r *contentModerationTestUserRepo) UpdateUserLastActiveAt(ctx context.Context, userID int64, activeAt time.Time) error {
	panic("unexpected UpdateUserLastActiveAt call")
placeholder

func (r *contentModerationTestUserRepo) UpdateBalance(ctx context.Context, id int64, amount float64) error {
	panic("unexpected UpdateBalance call")
placeholder

func (r *contentModerationTestUserRepo) DeductBalance(ctx context.Context, id int64, amount float64) error {
	panic("unexpected DeductBalance call")
placeholder

func (r *contentModerationTestUserRepo) UpdateConcurrency(ctx context.Context, id int64, amount int) error {
	panic("unexpected UpdateConcurrency call")
placeholder

func (r *contentModerationTestUserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	panic("unexpected ExistsByEmail call")
placeholder

func (r *contentModerationTestUserRepo) RemoveGroupFromAllowedGroups(ctx context.Context, groupID int64) (int64, error) {
	panic("unexpected RemoveGroupFromAllowedGroups call")
placeholder

func (r *contentModerationTestUserRepo) AddGroupToAllowedGroups(ctx context.Context, userID int64, groupID int64) error {
	panic("unexpected AddGroupToAllowedGroups call")
placeholder

func (r *contentModerationTestUserRepo) RemoveGroupFromUserAllowedGroups(ctx context.Context, userID int64, groupID int64) error {
	panic("unexpected RemoveGroupFromUserAllowedGroups call")
placeholder

func (r *contentModerationTestUserRepo) ListUserAuthIdentities(ctx context.Context, userID int64) ([]UserAuthIdentityRecord, error) {
	panic("unexpected ListUserAuthIdentities call")
placeholder

func (r *contentModerationTestUserRepo) UnbindUserAuthProvider(ctx context.Context, userID int64, provider string) error {
	panic("unexpected UnbindUserAuthProvider call")
placeholder

func (r *contentModerationTestUserRepo) UpdateTotpSecret(ctx context.Context, userID int64, encryptedSecret *string) error {
	panic("unexpected UpdateTotpSecret call")
placeholder

func (r *contentModerationTestUserRepo) EnableTotp(ctx context.Context, userID int64) error {
	panic("unexpected EnableTotp call")
placeholder

func (r *contentModerationTestUserRepo) DisableTotp(ctx context.Context, userID int64) error {
	panic("unexpected DisableTotp call")
placeholder

type contentModerationTestAuthCacheInvalidator struct {
	userIDs []int64
placeholder

func (i *contentModerationTestAuthCacheInvalidator) InvalidateAuthCacheByKey(ctx context.Context, key string) {
placeholder

func (i *contentModerationTestAuthCacheInvalidator) InvalidateAuthCacheByUserID(ctx context.Context, userID int64) {
	i.userIDs = append(i.userIDs, userID)
placeholder

func (i *contentModerationTestAuthCacheInvalidator) InvalidateAuthCacheByGroupID(ctx context.Context, groupID int64) {
placeholder

func (c *contentModerationTestHashCache) RecordFlaggedInputHash(ctx context.Context, inputHash string) error {
	if c.hashes == nil {
		c.hashes = map[string]struct{placeholder{placeholder
placeholder
	c.hashes[inputHash] = struct{placeholder{placeholder
	c.recorded = append(c.recorded, inputHash)
	return nil
placeholder

func (c *contentModerationTestHashCache) HasFlaggedInputHash(ctx context.Context, inputHash string) (bool, error) {
	c.checked = append(c.checked, inputHash)
	if c.hasResultUsed {
		return c.hasResult, nil
placeholder
	_, ok := c.hashes[inputHash]
	return ok, nil
placeholder

func (c *contentModerationTestHashCache) DeleteFlaggedInputHash(ctx context.Context, inputHash string) (bool, error) {
	c.deleted = append(c.deleted, inputHash)
	if c.hashes == nil {
		return false, nil
placeholder
	if _, ok := c.hashes[inputHash]; !ok {
		return false, nil
placeholder
	delete(c.hashes, inputHash)
	return true, nil
placeholder

func (c *contentModerationTestHashCache) ClearFlaggedInputHashes(ctx context.Context) (int64, error) {
	deleted := int64(len(c.hashes))
	c.hashes = map[string]struct{placeholder{placeholder
	return deleted, nil
placeholder

func (c *contentModerationTestHashCache) CountFlaggedInputHashes(ctx context.Context) (int64, error) {
	return int64(len(c.hashes)), nil
placeholder

func TestBuildContentModerationLog_RedactsInputExcerpt(t *testing.T) {
	svc := &ContentModerationService{placeholder
	cfg := defaultContentModerationConfig()
	input := ContentModerationCheckInput{
		RequestID: "req-1",
		Endpoint:  "/v1/chat/completions",
		Provider:  "openai",
placeholder

	log := svc.buildLog(input, cfg, ContentModerationActionAllow, true, "sexual", 0.8, map[string]float64{"sexual": 0.8placeholder, "hello sk-proj-1234567890abcdef", nil, nil, "")

	require.NotContains(t, log.InputExcerpt, "sk-proj-1234567890abcdef")
	require.Contains(t, log.InputExcerpt, "[已脱敏]")
placeholder

func TestRedactContentModerationSecrets_LongHexAndTokens(t *testing.T) {
	input := "你哈市多大事cf5bbdc4cd508f3aaf0d2070d529d4a4ac29099f8ecc357f696df28e1df91554 token=placeholder Bearer placeholder"

	out := redactContentModerationSecrets(input)

	require.NotContains(t, out, "cf5bbdc4cd508f3aaf0d2070d529d4a4ac29099f8ecc357f696df28e1df91554")
	require.NotContains(t, out, "placeholder")
	require.NotContains(t, out, "eyJhbGciOiJIUzI1NiJ9")
	require.Contains(t, out, "[已脱敏]")
placeholder

func TestContentModerationConfigNormalize_NonHitRetentionMaxThreeDays(t *testing.T) {
	cfg := defaultContentModerationConfig()
	cfg.NonHitRetentionDays = 30

	cfg.normalize()

	require.Equal(t, 3, cfg.NonHitRetentionDays)
placeholder

func TestExtractContentModerationInput_AnthropicImageSourceOnlyParticipatesInMemory(t *testing.T) {
	body := []byte(`{
		"messages": [
			{"role":"user","content":"old"placeholder,
			{"role":"assistant","content":"ok"placeholder,
			{"role":"user","content":[
				{"type":"text","text":"检查这张图"placeholder,
				{"type":"image","source":{"type":"base64","media_type":"image/png","data":"aGVsbG8="placeholderplaceholder
			]placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolAnthropicMessages, body)
	require.Equal(t, "检查这张图", input.Text)
	require.Equal(t, []string{"data:image/png;base64,aGVsbG8="placeholder, input.Images)

	log := (&ContentModerationService{placeholder).buildLog(ContentModerationCheckInput{placeholder, defaultContentModerationConfig(), ContentModerationActionAllow, false, "", 0, nil, input.ExcerptText(), nil, nil, "")
	require.Equal(t, "检查这张图", log.InputExcerpt)
	require.NotContains(t, log.InputExcerpt, "aGVsbG8=")
placeholder

func TestExtractContentModerationInput_AnthropicKeepsEphemeralUserTextAndSkipsSystemReminders(t *testing.T) {
	body := []byte(`{
		"messages": [
			{
				"role": "user",
				"content": [
					{"type": "text", "text": "<system-reminder>工具说明</system-reminder>"placeholder,
					{"type": "text", "text": "<system-reminder>Ainder>\n\n"placeholder,
					{"type": "text", "text": "hid", "cache_control": {"type": "ephemeral"placeholderplaceholder
				]
		placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolAnthropicMessages, body)

	require.Equal(t, "hid", input.Text)
	require.Empty(t, input.Images)
placeholder

func TestExtractContentModerationInput_OpenAIChatUsesLastUserMessage(t *testing.T) {
	body := []byte(`{
		"model":"gpt-5.5",
		"messages":[
			{"role":"system","content":"system prompt"placeholder,
			{"role":"user","content":"old user"placeholder,
			{"role":"assistant","content":"ok"placeholder,
			{"role":"user","content":[{"type":"text","text":"latest user"placeholder,{"type":"image_url","image_url":{"url":"https://example.com/a.png"placeholderplaceholder]placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolOpenAIChat, body)

	require.Equal(t, "latest user", input.Text)
	require.Equal(t, []string{"https://example.com/a.png"placeholder, input.Images)
	require.NotContains(t, input.Text, "old user")
	require.NotContains(t, input.Text, "system prompt")
placeholder

func TestExtractContentModerationInput_OpenAIImagesIncludesPromptAndImages(t *testing.T) {
	body := []byte(`{
		"prompt":"replace background",
		"images":[
			{"image_url":"https://example.com/source.png"placeholder,
			{"image_url":"data:image/png;base64,aGVsbG8="placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolOpenAIImages, body)

	require.Equal(t, "replace background", input.Text)
	require.Equal(t, []string{"https://example.com/source.png", "data:image/png;base64,aGVsbG8="placeholder, input.Images)
placeholder

func TestExtractContentModerationInput_OpenAIResponsesCodexPayloadUsesLastUserMessage(t *testing.T) {
	body := []byte(`{
		"model":"gpt-5.5",
		"instructions":"instructions.....",
		"input":[
			{"type":"message","role":"developer","content":[{"type":"input_text","text":"developer permissions sk-proj-1234567890abcdef"placeholder]placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"first user prompt"placeholder]placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"last user prompt"placeholder]placeholder
		],
		"prompt_cache_key":"cache-key"
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolOpenAIResponses, body)

	require.Equal(t, "last user prompt", input.Text)
	require.Empty(t, input.Images)
	require.NotContains(t, input.Text, "developer permissions")
	require.NotContains(t, input.Text, "first user prompt")
placeholder

func TestContentModerationCheck_OpenAIResponsesRecordsNonHitForCodexPayload(t *testing.T) {
	var moderationRequest moderationAPIRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/moderations", r.URL.Path)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&moderationRequest))
		_ = json.NewEncoder(w).Encode(moderationAPIResponse{
			Results: []moderationAPIResult{{
				CategoryScores: map[string]float64{"sexual": 0.01placeholder,
	placeholder
	placeholder)
placeholder))
	defer server.Close()

	cfg := defaultContentModerationConfig()
	cfg.Enabled = true
	cfg.Mode = ContentModerationModePreBlock
	cfg.BaseURL = server.URL
	cfg.APIKeys = []string{"sk-test"placeholder
	cfg.RecordNonHits = true
	rawCfg, err := json.Marshal(cfg)
placeholder

	repo := &contentModerationTestRepo{placeholder
	svc := NewContentModerationService(
		&contentModerationTestSettingRepo{values: map[string]string{
			SettingKeyRiskControlEnabled:      "true",
			SettingKeyContentModerationConfig: string(rawCfg),
placeholder
		repo,
		&contentModerationTestHashCache{placeholder,
		nil,
		nil,
		nil,
		nil,
	)

	body := []byte(`{
		"model":"gpt-5.5",
		"input":[
			{"type":"message","role":"developer","content":[{"type":"input_text","text":"developer instructions should not be audited"placeholder]placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"first user prompt"placeholder]placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"last user prompt"placeholder]placeholder
		]
placeholder`)
	decision, err := svc.Check(context.Background(), ContentModerationCheckInput{
		UserID:   1001,
		Endpoint: "/responses",
		Provider: "openai",
		Model:    "gpt-5.5",
		Protocol: ContentModerationProtocolOpenAIResponses,
		Body:     body,
placeholder)

placeholder
	require.False(t, decision.Blocked)
	require.Len(t, repo.logs, 1)
	require.False(t, repo.logs[0].Flagged)
	require.Equal(t, ContentModerationActionAllow, repo.logs[0].Action)
	require.Equal(t, "/responses", repo.logs[0].Endpoint)
	require.Equal(t, "last user prompt", repo.logs[0].InputExcerpt)
	require.Equal(t, "last user prompt", moderationRequest.Input)
placeholder

func TestContentModerationCheck_PreBlockBlocksCodexResponsesLatestUserInput(t *testing.T) {
	var moderationRequest moderationAPIRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/moderations", r.URL.Path)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&moderationRequest))
		_ = json.NewEncoder(w).Encode(moderationAPIResponse{
			Results: []moderationAPIResult{{
				CategoryScores: map[string]float64{"sexual": 0.9placeholder,
	placeholder
	placeholder)
placeholder))
	defer server.Close()

	cfg := defaultContentModerationConfig()
	cfg.Enabled = true
	cfg.Mode = ContentModerationModePreBlock
	cfg.BaseURL = server.URL
	cfg.APIKeys = []string{"sk-test"placeholder
	cfg.BlockStatus = http.StatusUnavailableForLegalReasons
	cfg.BlockMessage = "内容审计测试阻断"
	rawCfg, err := json.Marshal(cfg)
placeholder

	repo := &contentModerationTestRepo{placeholder
	svc := NewContentModerationService(
		&contentModerationTestSettingRepo{values: map[string]string{
			SettingKeyRiskControlEnabled:      "true",
			SettingKeyContentModerationConfig: string(rawCfg),
placeholder
		repo,
		&contentModerationTestHashCache{placeholder,
		nil,
		nil,
		nil,
		nil,
	)

	body := []byte(`{
		"model":"gpt-5.5",
		"instructions":"instructions.....",
		"input":[
			{"type":"message","role":"developer","content":[{"type":"input_text","text":"developer instructions should not be audited"placeholder]placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"environment context"placeholder]placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"latest blocked prompt"placeholder]placeholder
		]
placeholder`)
	decision, err := svc.Check(context.Background(), ContentModerationCheckInput{
		UserID:   1001,
		Endpoint: "/responses",
		Provider: "openai",
		Model:    "gpt-5.5",
		Protocol: ContentModerationProtocolOpenAIResponses,
		Body:     body,
placeholder)

placeholder
	require.True(t, decision.Blocked)
	require.Equal(t, ContentModerationActionBlock, decision.Action)
	require.Equal(t, http.StatusUnavailableForLegalReasons, decision.StatusCode)
	require.Equal(t, "内容审计测试阻断", decision.Message)
	require.Len(t, repo.logs, 1)
	require.True(t, repo.logs[0].Flagged)
	require.Equal(t, ContentModerationActionBlock, repo.logs[0].Action)
	require.Equal(t, ContentModerationModePreBlock, repo.logs[0].Mode)
	require.Equal(t, "latest blocked prompt", repo.logs[0].InputExcerpt)
	require.Equal(t, "latest blocked prompt", moderationRequest.Input)
placeholder

func TestBuildContentModerationTestAuditResult_UsesConfiguredThresholdsOnly(t *testing.T) {
	result := buildContentModerationTestAuditResult(&moderationAPIResult{
		Flagged: true,
		CategoryScores: map[string]float64{
			"harassment": 0.65,
	placeholder,
placeholder, nil)

	require.NotNil(t, result)
	require.False(t, result.Flagged)
	require.Equal(t, "harassment", result.HighestCategory)
	require.Equal(t, 0.65, result.HighestScore)
	require.Equal(t, 0.65, result.CompositeScore)
	require.Equal(t, 0.98, result.Thresholds["harassment"])
placeholder

func TestContentModerationCheck_PreHashUsesRedisHashCache(t *testing.T) {
	cfg := defaultContentModerationConfig()
	cfg.Enabled = true
	cfg.PreHashCheckEnabled = true
	cfg.APIKeys = []string{"sk-test"placeholder
	cfg.BlockStatus = http.StatusConflict
	cfg.BlockMessage = "命中历史风险输入"
	rawCfg, err := json.Marshal(cfg)
placeholder

	hashCache := &contentModerationTestHashCache{hashes: map[string]struct{placeholder{placeholderplaceholder
	content := ContentModerationInput{Text: "blocked prompt"placeholder
	content.Normalize()
	hashCache.hashes[content.Hash()] = struct{placeholder{placeholder

	svc := NewContentModerationService(
		&contentModerationTestSettingRepo{values: map[string]string{
			SettingKeyRiskControlEnabled:      "true",
			SettingKeyContentModerationConfig: string(rawCfg),
placeholder
		&contentModerationTestRepo{placeholder,
		hashCache,
		nil,
		nil,
		nil,
		nil,
	)

	decision, err := svc.Check(context.Background(), ContentModerationCheckInput{
		Protocol: ContentModerationProtocolOpenAIChat,
		Body:     []byte(`{"messages":[{"role":"user","content":"blocked prompt"placeholder]placeholder`),
placeholder)
placeholder
	require.True(t, decision.Blocked)
	require.Equal(t, ContentModerationActionHashBlock, decision.Action)
	require.Equal(t, http.StatusConflict, decision.StatusCode)
	require.Equal(t, content.Hash(), decision.InputHash)
	require.Contains(t, decision.Message, "命中历史风险输入")
	require.Contains(t, decision.Message, content.Hash())
	require.Len(t, hashCache.checked, 1)
placeholder

func TestContentModerationCheck_PreBlockFlaggedWritesRedisHashCache(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		_ = json.NewEncoder(w).Encode(moderationAPIResponse{
			Results: []moderationAPIResult{{
				CategoryScores: map[string]float64{"sexual": 0.9placeholder,
	placeholder
	placeholder)
placeholder))
	defer server.Close()

	cfg := defaultContentModerationConfig()
	cfg.Enabled = true
	cfg.Mode = ContentModerationModePreBlock
	cfg.PreHashCheckEnabled = true
	cfg.BaseURL = server.URL
	cfg.APIKeys = []string{"sk-test"placeholder
	cfg.BlockStatus = http.StatusConflict
	cfg.BlockMessage = "命中风险输入"
	rawCfg, err := json.Marshal(cfg)
placeholder

	repo := &contentModerationTestRepo{placeholder
	hashCache := &contentModerationTestHashCache{placeholder
	svc := NewContentModerationService(
		&contentModerationTestSettingRepo{values: map[string]string{
			SettingKeyRiskControlEnabled:      "true",
			SettingKeyContentModerationConfig: string(rawCfg),
placeholder
		repo,
		hashCache,
		nil,
		nil,
		nil,
		nil,
	)

	body := []byte(`{"messages":[{"role":"user","content":"repeat blocked prompt"placeholder]placeholder`)
	decision, err := svc.Check(context.Background(), ContentModerationCheckInput{
		Protocol: ContentModerationProtocolOpenAIChat,
		Body:     body,
placeholder)
placeholder
	require.True(t, decision.Blocked)
	require.Equal(t, ContentModerationActionBlock, decision.Action)
	require.Equal(t, 1, requestCount)
	require.Len(t, hashCache.recorded, 1)
	require.Len(t, repo.logs, 1)

	decision, err = svc.Check(context.Background(), ContentModerationCheckInput{
		Protocol: ContentModerationProtocolOpenAIChat,
		Body:     body,
placeholder)
placeholder
	require.True(t, decision.Blocked)
	require.Equal(t, ContentModerationActionHashBlock, decision.Action)
	require.Equal(t, hashCache.recorded[0], decision.InputHash)
	require.Equal(t, 1, requestCount)
	require.Len(t, repo.logs, 1)
placeholder

func TestContentModerationDeleteFlaggedInputHash_NormalizesAndDeletes(t *testing.T) {
	existingHash := strings.Repeat("a", 64)
	hashCache := &contentModerationTestHashCache{hashes: map[string]struct{placeholder{
		existingHash: {placeholder,
placeholderplaceholder
	svc := &ContentModerationService{hashCache: hashCacheplaceholder

	result, err := svc.DeleteFlaggedInputHash(context.Background(), strings.ToUpper(existingHash))

placeholder
	require.Equal(t, existingHash, result.InputHash)
	require.True(t, result.Deleted)
	require.NotContains(t, hashCache.hashes, existingHash)
	require.Equal(t, []string{existingHashplaceholder, hashCache.deleted)

	result, err = svc.DeleteFlaggedInputHash(context.Background(), existingHash)

placeholder
	require.Equal(t, existingHash, result.InputHash)
	require.False(t, result.Deleted)
placeholder

func TestContentModerationClearFlaggedInputHashesAndStatusCount(t *testing.T) {
	cfg := defaultContentModerationConfig()
	cfg.Enabled = true
	rawCfg, err := json.Marshal(cfg)
placeholder

	hashCache := &contentModerationTestHashCache{hashes: map[string]struct{placeholder{
		strings.Repeat("a", 64): {placeholder,
		strings.Repeat("b", 64): {placeholder,
placeholderplaceholder
	svc := &ContentModerationService{
		settingRepo: &contentModerationTestSettingRepo{values: map[string]string{
			SettingKeyRiskControlEnabled:      "true",
			SettingKeyContentModerationConfig: string(rawCfg),
placeholder
		hashCache: hashCache,
		keyHealth: make(map[string]*contentModerationKeyHealth),
placeholder

	status, err := svc.GetStatus(context.Background())
placeholder
	require.Equal(t, int64(2), status.FlaggedHashCount)

	result, err := svc.ClearFlaggedInputHashes(context.Background())
placeholder
	require.Equal(t, int64(2), result.Deleted)

	status, err = svc.GetStatus(context.Background())
placeholder
	require.Equal(t, int64(0), status.FlaggedHashCount)
placeholder

func TestContentModerationCheck_AsyncFlaggedWritesRedisHashCache(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(moderationAPIResponse{
			Results: []moderationAPIResult{{
				CategoryScores: map[string]float64{"sexual": 0.9placeholder,
	placeholder
	placeholder)
placeholder))
	defer server.Close()

	cfg := defaultContentModerationConfig()
	cfg.Enabled = true
	cfg.Mode = ContentModerationModeObserve
	cfg.BaseURL = server.URL
	cfg.APIKeys = []string{"sk-test"placeholder
	rawCfg, err := json.Marshal(cfg)
placeholder

	repo := &contentModerationTestRepo{placeholder
	hashCache := &contentModerationTestHashCache{placeholder
	svc := NewContentModerationService(
		&contentModerationTestSettingRepo{values: map[string]string{
			SettingKeyRiskControlEnabled:      "true",
			SettingKeyContentModerationConfig: string(rawCfg),
placeholder
		repo,
		hashCache,
		nil,
		nil,
		nil,
		nil,
	)

	decision := svc.checkSync(context.Background(), ContentModerationCheckInput{
		Protocol: ContentModerationProtocolOpenAIChat,
		Body:     []byte(`{"messages":[{"role":"user","content":"bad prompt"placeholder]placeholder`),
placeholder, cfg, ContentModerationInput{Text: "bad prompt"placeholder, strings.Repeat("b", 64), contentModerationIntPtr(25), false)

	require.False(t, decision.Blocked)
	require.Len(t, hashCache.recorded, 1)
	require.Len(t, repo.logs, 1)
placeholder

func TestBuildContentModerationAccountDisabledEmailBody_ContainsBanDetails(t *testing.T) {
	userID := int64(1001)
	cfg := defaultContentModerationConfig()
	cfg.BanThreshold = 10
	body := buildContentModerationAccountDisabledEmailBody("Sub2API <Admin>", &ContentModerationLog{
		UserID:          &userID,
		UserEmail:       "user@example.com",
		GroupName:       "vip_2",
		HighestCategory: "sexual",
		HighestScore:    0.926,
		ViolationCount:  10,
placeholder, cfg)

	require.Contains(t, body, "账户已被自动禁用")
	require.Contains(t, body, "封禁详情")
	require.Contains(t, body, "账户当前处于封禁状态，所有 API 请求将被拒绝")
	require.Contains(t, body, "10 次（阈值 10）")
	require.Contains(t, body, "sexual / 0.926")
	require.Contains(t, body, "Sub2API &lt;Admin&gt;")
placeholder

func TestContentModerationUnbanUser_ActivatesUserAndInvalidatesAuthCache(t *testing.T) {
	userRepo := &contentModerationTestUserRepo{user: &User{ID: 1001, Email: "user@example.com", Status: StatusDisabledplaceholderplaceholder
	invalidator := &contentModerationTestAuthCacheInvalidator{placeholder
	repo := &contentModerationTestRepo{placeholder
	svc := NewContentModerationService(nil, repo, nil, nil, userRepo, invalidator, nil)

	result, err := svc.UnbanUser(context.Background(), 1001)

placeholder
	require.Equal(t, int64(1001), result.UserID)
	require.Equal(t, StatusActive, result.Status)
	require.Len(t, userRepo.updated, 1)
	require.Equal(t, StatusActive, userRepo.updated[0].Status)
	require.Equal(t, []int64{1001placeholder, invalidator.userIDs)
placeholder

func TestContentModerationUnbanUser_ActiveUserOnlyInvalidatesAuthCache(t *testing.T) {
	userRepo := &contentModerationTestUserRepo{user: &User{ID: 1001, Email: "user@example.com", Status: StatusActiveplaceholderplaceholder
	invalidator := &contentModerationTestAuthCacheInvalidator{placeholder
	repo := &contentModerationTestRepo{placeholder
	svc := NewContentModerationService(nil, repo, nil, nil, userRepo, invalidator, nil)

	result, err := svc.UnbanUser(context.Background(), 1001)

placeholder
	require.Equal(t, StatusActive, result.Status)
	require.Empty(t, userRepo.updated)
	require.Equal(t, []int64{1001placeholder, invalidator.userIDs)
placeholder

func contentModerationIntPtr(v int) *int {
	return &v
placeholder
