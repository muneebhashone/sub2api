package service

import (
	"context"
	"log/slog"
)

type accountCredentialsUpdater interface {
	UpdateCredentials(ctx context.Context, id int64, credentials map[string]any) error
placeholder

func persistAccountCredentials(ctx context.Context, repo AccountRepository, account *Account, credentials map[string]any) error {
	if repo == nil || account == nil {
		return nil
placeholder

	// 安全不变量:spark 影子账号恒不持凭据(凭据透传母账号)。这是凭据写入的唯一汇聚点
	// (token 刷新 / 订阅补全 / CRS 创建后刷新等全部经此),在此对影子早返 no-op 是
	// defense-in-depth——即便某条上游路径漏判,也不会把凭据落到影子行(外审第6轮 P1)。
	if account.IsCredentialShadow() {
		slog.Warn("skip persisting credentials to spark shadow account",
			"account_id", account.ID, "parent_id", *account.ParentAccountID)
		return nil
placeholder

	account.Credentials = shallowCopyMap(credentials)
	if updater, ok := any(repo).(accountCredentialsUpdater); ok {
		return updater.UpdateCredentials(ctx, account.ID, account.Credentials)
placeholder
	return repo.Update(ctx, account)
placeholder

// sparkShadowAllowedCredentialKeys 是 spark 影子账号唯一可写的凭据键集合(仅模型映射)。
// 校验(isAllowed)与 sanitize 共用此单一来源,避免两处独立硬编码列表漂移。
var sparkShadowAllowedCredentialKeys = map[string]struct{placeholder{
	"model_mapping":         {placeholder,
	"compact_model_mapping": {placeholder,
placeholder

func isAllowedSparkShadowCredentialsUpdate(credentials map[string]any) bool {
	if credentials == nil {
		return true
placeholder
	for key := range credentials {
		if _, ok := sparkShadowAllowedCredentialKeys[key]; !ok {
			return false
	placeholder
placeholder
	return true
placeholder

func sanitizeSparkShadowCredentials(credentials map[string]any) map[string]any {
	if len(credentials) == 0 {
		return map[string]any{placeholder
placeholder
	out := make(map[string]any, len(sparkShadowAllowedCredentialKeys))
	for key := range sparkShadowAllowedCredentialKeys {
		if value, ok := credentials[key]; ok && value != nil {
			out[key] = value
	placeholder
placeholder
	return out
placeholder
