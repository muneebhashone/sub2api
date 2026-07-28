//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// api_keys 的 quota_used / usage_5h|1d|7d 由计费热路径原子递增。
// 编辑 Key（改名、换分组……）若整行回写，并发累计的用量就会被旧快照覆盖。
// 这些用例锁死"只声明请求真正要改的列"。

type updateFieldsAPIKeyRepoStub struct {
	quotaBaseAPIKeyRepoStub
	key          *APIKey
	updateFields []APIKeyUpdateFields
placeholder

// IncrementQuotaUsed 模拟计费热路径上的原子递增：只动 quota_used。
func (s *updateFieldsAPIKeyRepoStub) IncrementQuotaUsed(_ context.Context, _ int64, amount float64) (float64, error) {
	s.key.QuotaUsed += amount
	return s.key.QuotaUsed, nil
placeholder

func (s *updateFieldsAPIKeyRepoStub) GetByID(context.Context, int64) (*APIKey, error) {
	clone := *s.key
	return &clone, nil
placeholder

func (s *updateFieldsAPIKeyRepoStub) Update(_ context.Context, _ *APIKey, fields APIKeyUpdateFields) error {
	s.updateFields = append(s.updateFields, fields)
	return nil
placeholder

func newUpdateFieldsAPIKeyService(key *APIKey) (*APIKeyService, *updateFieldsAPIKeyRepoStub) {
	repo := &updateFieldsAPIKeyRepoStub{key: keyplaceholder
	return &APIKeyService{apiKeyRepo: repoplaceholder, repo
placeholder

func TestAPIKeyUpdate_OnlyDeclaresRequestedColumns(t *testing.T) {
	name := "renamed"
	quota := 500.0
	rateLimit := 42.0
	whitelist := []string{"10.0.0.1"placeholder

	tests := []struct {
		name string
		req  UpdateAPIKeyRequest
		want APIKeyUpdateFields
placeholder{
		{
			name: "name only",
			req:  UpdateAPIKeyRequest{Name: &nameplaceholder,
			want: APIKeyUpdateFields{Name: trueplaceholder,
	placeholder,
		{
			name: "quota only",
			req:  UpdateAPIKeyRequest{Quota: &quotaplaceholder,
			want: APIKeyUpdateFields{Quota: trueplaceholder,
	placeholder,
		{
			name: "rate limit threshold only",
			req:  UpdateAPIKeyRequest{RateLimit5h: &rateLimitplaceholder,
			want: APIKeyUpdateFields{RateLimits: trueplaceholder,
	placeholder,
		{
			name: "ip whitelist only",
			req:  UpdateAPIKeyRequest{IPWhitelist: &whitelistplaceholder,
			want: APIKeyUpdateFields{IPRules: trueplaceholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, repo := newUpdateFieldsAPIKeyService(&APIKey{
				ID:        1,
				UserID:    7,
				Key:       "sk-test",
				Name:      "before",
				Status:    StatusActive,
				Quota:     100,
				QuotaUsed: 30,
				Usage5h:   12,
		placeholder)

			_, err := svc.Update(context.Background(), 1, 7, tt.req)
		placeholder
			require.Equal(t, []APIKeyUpdateFields{tt.wantplaceholder, repo.updateFields)
	placeholder)
placeholder
placeholder

// 显式重置仍需声明对应的列，避免收窄写入列时把功能改坏。
func TestAPIKeyUpdate_DeclaresUsageColumnsOnExplicitReset(t *testing.T) {
	reset := true
	svc, repo := newUpdateFieldsAPIKeyService(&APIKey{
		ID: 1, UserID: 7, Key: "sk-test", Status: StatusActive, Quota: 100, QuotaUsed: 30, Usage5h: 12,
placeholder)

	_, err := svc.Update(context.Background(), 1, 7, UpdateAPIKeyRequest{
		ResetQuota:          &reset,
		ResetRateLimitUsage: &reset,
placeholder)
placeholder
	require.Equal(t, []APIKeyUpdateFields{{QuotaUsed: true, RateLimitUsage: trueplaceholderplaceholder, repo.updateFields)
placeholder

// 配额扩容会顺带把 quota_exhausted 复活为 active，此时必须声明 status。
func TestAPIKeyUpdate_DeclaresStatusWhenReactivated(t *testing.T) {
	quota := 500.0
	svc, repo := newUpdateFieldsAPIKeyService(&APIKey{
		ID: 1, UserID: 7, Key: "sk-test", Status: StatusAPIKeyQuotaExhausted, Quota: 100, QuotaUsed: 100,
placeholder)

	_, err := svc.Update(context.Background(), 1, 7, UpdateAPIKeyRequest{Quota: &quotaplaceholder)
placeholder
	require.Equal(t, []APIKeyUpdateFields{{Quota: true, Status: trueplaceholderplaceholder, repo.updateFields)
placeholder

// 计费热路径把 Key 标记为配额耗尽时只写 status，
// 否则会把刚原子递增的 quota_used 按快照覆盖掉。
func TestUpdateQuotaUsed_ExhaustedMarkOnlyDeclaresStatus(t *testing.T) {
	repo := &updateFieldsAPIKeyRepoStub{key: &APIKey{
		ID: 1, UserID: 7, Key: "sk-test", Status: StatusActive, Quota: 10, QuotaUsed: 10,
placeholderplaceholder
	svc := &APIKeyService{apiKeyRepo: repoplaceholder

	require.NoError(t, svc.UpdateQuotaUsed(context.Background(), 1, 5))
	require.Equal(t, []APIKeyUpdateFields{{Status: trueplaceholderplaceholder, repo.updateFields)
placeholder
