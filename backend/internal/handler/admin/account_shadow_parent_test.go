package admin

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestEnrichShadowParentInfo(t *testing.T) {
	pid := int64(100)
	parent := &service.Account{
		ID: 100,
placeholder
			"email":                   "owner@example.com",
			"plan_type":               "pro",
			"subscription_expires_at": "2026-12-31T00:00:00Z",
			"chatgpt_account_id":      "acct_123",
	placeholder,
		Extra: map[string]any{"privacy_mode": "training_off"placeholder,
placeholder
	parents := map[int64]*service.Account{100: parentplaceholder

	shadow := AccountWithConcurrency{Account: &dto.Account{ID: 200, ParentAccountID: &pidplaceholderplaceholder
	normal := AccountWithConcurrency{Account: &dto.Account{ID: 1placeholderplaceholder
	orphan := AccountWithConcurrency{Account: &dto.Account{ID: 201, ParentAccountID: ptrInt64(999)placeholderplaceholder
	items := []AccountWithConcurrency{shadow, normal, orphanplaceholder

	enrichShadowParentInfo(items, parents)

	require.Equal(t, "owner@example.com", items[0].ParentEmail, "影子回填母账号邮箱")
	require.Equal(t, "pro", items[0].ParentPlanType)
	require.Equal(t, "training_off", items[0].ParentPrivacyMode)
	require.Equal(t, "2026-12-31T00:00:00Z", items[0].ParentSubscriptionExpiresAt)
	require.Equal(t, "acct_123", items[0].ParentChatGPTAccountID)

	require.Empty(t, items[1].ParentEmail, "非影子不回填")
	require.Empty(t, items[2].ParentEmail, "母账号缺失时优雅留空")
placeholder

func ptrInt64(v int64) *int64 { return &v placeholder
