package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGetAccessToken_SparkShadowResolvesToParent 验证对影子账号调用 GetAccessToken
// 时能透明地解析到母账号的凭据，防止 refresh_token 脱钩。
// 影子账号不持凭据；断言必须返回母账号的 access_token。
func TestGetAccessToken_SparkShadowResolvesToParent(t *testing.T) {
	ctx := context.Background()

	parentID := int64(100)
	parent := Account{
		ID:       parentID,
		Platform: PlatformOpenAI,
		Type:     AccountTypeOAuth,
		Status:   StatusActive,
placeholder
			"access_token": "parent-access-token",
	placeholder,
placeholder
	shadow := Account{
		ID:              200,
		Platform:        PlatformOpenAI,
		Type:            AccountTypeOAuth,
		ParentAccountID: &parentID,
		// 影子账号不持凭据，与生产语义一致
placeholder

	repo := &stubOpenAIAccountRepo{accounts: []Account{parentplaceholderplaceholder

	svc := &OpenAIGatewayService{
		accountRepo: repo,
		// openAITokenProvider=nil → 走降级路径，直接读 account.GetOpenAIAccessToken()
placeholder

	// Before fix (RED): shadow 无凭据 → GetOpenAIAccessToken()="" → error
	// After fix (GREEN): shadow 解析到 parent → 返回 parent 的 "parent-access-token"
	token, tokenType, err := svc.GetAccessToken(ctx, &shadow)
placeholder
	require.Equal(t, "parent-access-token", token)
	require.Equal(t, "oauth", tokenType)
placeholder
