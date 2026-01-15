package service

import "strconv"

// OpenAITokenCacheKey 生成 OpenAI OAuth 账号的缓存键
// 格式: "openai:account:{account_idplaceholder"
func OpenAITokenCacheKey(account *Account) string {
	return "openai:account:" + strconv.FormatInt(account.ID, 10)
placeholder

// ClaudeTokenCacheKey 生成 Claude (Anthropic) OAuth 账号的缓存键
// 格式: "claude:account:{account_idplaceholder"
func ClaudeTokenCacheKey(account *Account) string {
	return "claude:account:" + strconv.FormatInt(account.ID, 10)
placeholder
