package service

import "github.com/gin-gonic/gin"

const errorPassthroughServiceContextKey = "error_passthrough_service"

// BindErrorPassthroughService 将错误透传服务绑定到请求上下文，供 service 层在非 failover 场景下复用规则。
func BindErrorPassthroughService(c *gin.Context, svc *ErrorPassthroughService) {
	if c == nil || svc == nil {
		return
placeholder
	c.Set(errorPassthroughServiceContextKey, svc)
placeholder

func getBoundErrorPassthroughService(c *gin.Context) *ErrorPassthroughService {
	if c == nil {
		return nil
placeholder
	v, ok := c.Get(errorPassthroughServiceContextKey)
	if !ok {
		return nil
placeholder
	svc, ok := v.(*ErrorPassthroughService)
	if !ok {
		return nil
placeholder
	return svc
placeholder

// applyErrorPassthroughRule 按规则改写错误响应；未命中时返回默认响应参数。
func applyErrorPassthroughRule(
	c *gin.Context,
	platform string,
	upstreamStatus int,
	responseBody []byte,
	defaultStatus int,
	defaultErrType string,
	defaultErrMsg string,
) (status int, errType string, errMsg string, matched bool) {
	status = defaultStatus
	errType = defaultErrType
	errMsg = defaultErrMsg

	svc := getBoundErrorPassthroughService(c)
	if svc == nil {
		return status, errType, errMsg, false
placeholder

	rule := svc.MatchRule(platform, upstreamStatus, responseBody)
	if rule == nil {
		return status, errType, errMsg, false
placeholder

	status = upstreamStatus
	if !rule.PassthroughCode && rule.ResponseCode != nil {
		status = *rule.ResponseCode
placeholder

	errMsg = ExtractUpstreamErrorMessage(responseBody)
	if !rule.PassthroughBody && rule.CustomMessage != nil {
		errMsg = *rule.CustomMessage
placeholder

	// 与现有 failover 场景保持一致：命中规则时统一返回 upstream_error。
	errType = "upstream_error"
	return status, errType, errMsg, true
placeholder
