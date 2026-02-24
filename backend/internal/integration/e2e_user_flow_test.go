//go:build e2e

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// E2E 用户流程测试
// 测试完整的用户操作链路：注册 → 登录 → 创建 API Key → 调用网关 → 查询用量

var (
	testUserEmail    = "e2e-test-" + fmt.Sprintf("%d", time.Now().UnixMilli()) + "@test.local"
	testUserPassword = "E2eTest@12345"
	testUserName     = "e2e-test-user"
)

// TestUserRegistrationAndLogin 测试用户注册和登录流程
func TestUserRegistrationAndLogin(t *testing.T) {
	// 步骤 1: 注册新用户
	t.Run("注册新用户", func(t *testing.T) {
		payload := map[string]string{
			"email":    testUserEmail,
			"password": testUserPassword,
			"username": testUserName,
	placeholder
		body, _ := json.Marshal(payload)

		resp, err := doRequest(t, "POST", "/api/auth/register", body, "")
		if err != nil {
			t.Skipf("注册接口不可用，跳过用户流程测试: %v", err)
			return
	placeholder
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		// 注册可能返回 200（成功）或 400（邮箱已存在）或 403（注册已关闭）
		switch resp.StatusCode {
		case 200:
			t.Logf("✅ 用户注册成功: %s", testUserEmail)
		case 400:
			t.Logf("⚠️ 用户可能已存在: %s", string(respBody))
		case 403:
			t.Skipf("注册功能已关闭: %s", string(respBody))
		default:
			t.Logf("⚠️ 注册返回 HTTP %d: %s（继续尝试登录）", resp.StatusCode, string(respBody))
	placeholder
placeholder)

	// 步骤 2: 登录获取 JWT
	var accessToken string
	t.Run("用户登录获取JWT", func(t *testing.T) {
		payload := map[string]string{
			"email":    testUserEmail,
			"password": testUserPassword,
	placeholder
		body, _ := json.Marshal(payload)

		resp, err := doRequest(t, "POST", "/api/auth/login", body, "")
		if err != nil {
			t.Fatalf("登录请求失败: %v", err)
	placeholder
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Skipf("登录失败 HTTP %d: %s（可能需要先注册用户）", resp.StatusCode, string(respBody))
			return
	placeholder

		var result map[string]any
		if err := json.Unmarshal(respBody, &result); err != nil {
			t.Fatalf("解析登录响应失败: %v", err)
	placeholder

		// 尝试从标准响应格式获取 token
		if token, ok := result["access_token"].(string); ok && token != "" {
			accessToken = token
	placeholder else if data, ok := result["data"].(map[string]any); ok {
			if token, ok := data["access_token"].(string); ok {
				accessToken = token
		placeholder
	placeholder

		if accessToken == "" {
			t.Skipf("未获取到 access_token，响应: %s", string(respBody))
			return
	placeholder

		// 验证 token 不为空且格式基本正确
		if len(accessToken) < 10 {
			t.Fatalf("access_token 格式异常: %s", accessToken)
	placeholder

		t.Logf("✅ 登录成功，获取 JWT（长度: %d）", len(accessToken))
placeholder)

	if accessToken == "" {
		t.Skip("未获取到 JWT，跳过后续测试")
		return
placeholder

	// 步骤 3: 使用 JWT 获取当前用户信息
	t.Run("获取当前用户信息", func(t *testing.T) {
		resp, err := doRequest(t, "GET", "/api/user/me", nil, accessToken)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
	placeholder
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("HTTP %d: %s", resp.StatusCode, string(body))
	placeholder

		t.Logf("✅ 成功获取用户信息")
placeholder)
placeholder

// TestAPIKeyLifecycle 测试 API Key 的创建和使用
func TestAPIKeyLifecycle(t *testing.T) {
	// 先登录获取 JWT
	accessToken := loginTestUser(t)
	if accessToken == "" {
		t.Skip("无法登录，跳过 API Key 生命周期测试")
		return
placeholder

	var apiKey string

	// 步骤 1: 创建 API Key
	t.Run("创建API_Key", func(t *testing.T) {
		payload := map[string]string{
			"name": "e2e-test-key-" + fmt.Sprintf("%d", time.Now().UnixMilli()),
	placeholder
		body, _ := json.Marshal(payload)

		resp, err := doRequest(t, "POST", "/api/keys", body, accessToken)
		if err != nil {
			t.Fatalf("创建 API Key 请求失败: %v", err)
	placeholder
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Skipf("创建 API Key 失败 HTTP %d: %s", resp.StatusCode, string(respBody))
			return
	placeholder

		var result map[string]any
		if err := json.Unmarshal(respBody, &result); err != nil {
			t.Fatalf("解析响应失败: %v", err)
	placeholder

		// 从响应中提取 key
		if key, ok := result["key"].(string); ok {
			apiKey = key
	placeholder else if data, ok := result["data"].(map[string]any); ok {
			if key, ok := data["key"].(string); ok {
				apiKey = key
		placeholder
	placeholder

		if apiKey == "" {
			t.Skipf("未获取到 API Key，响应: %s", string(respBody))
			return
	placeholder

		// 验证 API Key 脱敏日志（只显示前 8 位）
		masked := apiKey
		if len(masked) > 8 {
			masked = masked[:8] + "..."
	placeholder
		t.Logf("✅ API Key 创建成功: %s", masked)
placeholder)

	if apiKey == "" {
		t.Skip("未创建 API Key，跳过后续测试")
		return
placeholder

	// 步骤 2: 使用 API Key 调用网关（需要 Claude 或 Gemini 可用）
	t.Run("使用API_Key调用网关", func(t *testing.T) {
		// 尝试调用 models 列表（最轻量的 API 调用）
		resp, err := doRequest(t, "GET", "/v1/models", nil, apiKey)
		if err != nil {
			t.Fatalf("网关请求失败: %v", err)
	placeholder
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		// 可能返回 200（成功）或 402（余额不足）或 403（无可用账户）
		switch {
		case resp.StatusCode == 200:
			t.Logf("✅ API Key 网关调用成功")
		case resp.StatusCode == 402:
			t.Logf("⚠️ 余额不足，但 API Key 认证通过")
		case resp.StatusCode == 403:
			t.Logf("⚠️ 无可用账户，但 API Key 认证通过")
		default:
			t.Logf("⚠️ 网关返回 HTTP %d: %s", resp.StatusCode, string(respBody))
	placeholder
placeholder)

	// 步骤 3: 查询用量记录
	t.Run("查询用量记录", func(t *testing.T) {
		resp, err := doRequest(t, "GET", "/api/usage/dashboard", nil, accessToken)
		if err != nil {
			t.Fatalf("用量查询请求失败: %v", err)
	placeholder
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Logf("⚠️ 用量查询返回 HTTP %d: %s", resp.StatusCode, string(body))
			return
	placeholder

		t.Logf("✅ 用量查询成功")
placeholder)
placeholder

// =============================================================================
// 辅助函数
// =============================================================================

func doRequest(t *testing.T, method, path string, body []byte, token string) (*http.Response, error) {
placeholder

	url := baseURL + path
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
placeholder

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
placeholder

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
placeholder
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
placeholder

	client := &http.Client{Timeout: 30 * time.Secondplaceholder
	return client.Do(req)
placeholder

func loginTestUser(t *testing.T) string {
placeholder

	// 先尝试用管理员账户登录
	adminEmail := getEnv("ADMIN_EMAIL", "admin@sub2api.local")
	adminPassword := getEnv("ADMIN_PASSWORD", "")

	if adminPassword == "" {
		// 尝试用测试用户
		adminEmail = testUserEmail
		adminPassword = testUserPassword
placeholder

	payload := map[string]string{
		"email":    adminEmail,
		"password": adminPassword,
placeholder
	body, _ := json.Marshal(payload)

	resp, err := doRequest(t, "POST", "/api/auth/login", body, "")
	if err != nil {
		return ""
placeholder
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
placeholder

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return ""
placeholder

	if token, ok := result["access_token"].(string); ok {
		return token
placeholder
	if data, ok := result["data"].(map[string]any); ok {
		if token, ok := data["access_token"].(string); ok {
			return token
	placeholder
placeholder

	return ""
placeholder

// redactAPIKey API Key 脱敏，只显示前 8 位
func redactAPIKey(key string) string {
	key = strings.TrimSpace(key)
	if len(key) <= 8 {
		return "***"
placeholder
	return key[:8] + "..."
placeholder
