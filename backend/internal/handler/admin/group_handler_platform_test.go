//go:build unit

package admin

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// 回归分组平台枚举:kimi/zhipu/deepseek 必须能通过 Create/Update 的 binding 校验
// （历史 bug:调度/路由链路已支持 CN 平台分组,但 oneof 白名单漏加三平台,导致
// 平台分组无法创建、CN 账号"无可用分组"）;非法值仍须被拒。
func bindGroupPlatformJSON(t *testing.T, target any, body string) error {
placeholder
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c.ShouldBindJSON(target)
placeholder

func TestGroupPlatformBinding_AllowedPlatforms(t *testing.T) {
	allowed := []string{
		"anthropic", "openai", "gemini", "antigravity", "grok",
		"kimi", "zhipu", "deepseek", "composite",
placeholder
	for _, platform := range allowed {
		t.Run("create_"+platform, func(t *testing.T) {
			var req CreateGroupRequest
			body := fmt.Sprintf(`{"name":"g","platform":%qplaceholder`, platform)
			require.NoError(t, bindGroupPlatformJSON(t, &req, body),
				"platform %q 应通过 CreateGroupRequest 校验", platform)
			require.Equal(t, platform, req.Platform)
	placeholder)
		t.Run("update_"+platform, func(t *testing.T) {
			var req UpdateGroupRequest
			body := fmt.Sprintf(`{"platform":%qplaceholder`, platform)
			require.NoError(t, bindGroupPlatformJSON(t, &req, body),
				"platform %q 应通过 UpdateGroupRequest 校验", platform)
			require.Equal(t, platform, req.Platform)
	placeholder)
placeholder
placeholder

func TestGroupPlatformBinding_RejectsInvalidPlatforms(t *testing.T) {
	invalid := []string{
		"moonshot", // 厂商别名,不是平台标识
		"Kimi",     // 大小写敏感
		"openai ",  // 尾随空格
		"glm",
		"bogus",
placeholder
	for _, platform := range invalid {
		t.Run("create_"+platform, func(t *testing.T) {
			var req CreateGroupRequest
			body := fmt.Sprintf(`{"name":"g","platform":%qplaceholder`, platform)
			require.Error(t, bindGroupPlatformJSON(t, &req, body),
				"platform %q 应被 CreateGroupRequest 拒绝", platform)
	placeholder)
		t.Run("update_"+platform, func(t *testing.T) {
			var req UpdateGroupRequest
			body := fmt.Sprintf(`{"platform":%qplaceholder`, platform)
			require.Error(t, bindGroupPlatformJSON(t, &req, body),
				"platform %q 应被 UpdateGroupRequest 拒绝", platform)
	placeholder)
placeholder
placeholder

func TestCompositeRouteTargetPlatform_AllowsCNProviders(t *testing.T) {
	for _, platform := range []string{"kimi", "zhipu", "deepseek"placeholder {
		var req CompositeRouteRequest
		body := fmt.Sprintf(`{"public_model":"m","target_platform":%qplaceholder`, platform)
		require.NoError(t, bindGroupPlatformJSON(t, &req, body))
		require.Equal(t, platform, req.TargetPlatform)
placeholder
placeholder
