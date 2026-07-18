package admin

import (
	"net/http"
	"net/netip"
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

var ingressRejectReasons = map[string]struct{placeholder{
	"query_api_key_deprecated": {placeholder, "api_key_required": {placeholder, "invalid_api_key": {placeholder,
	"invalid_auth_rate_limited": {placeholder,
	"api_key_auth_overloaded":   {placeholder,
	"api_key_disabled":          {placeholder, "ip_restricted": {placeholder, "user_inactive": {placeholder, "group_deleted": {placeholder,
	"group_disabled": {placeholder, "group_not_allowed": {placeholder, "group_unassigned": {placeholder, "other": {placeholder,
placeholder

var ingressRejectRouteFamilies = map[string]struct{placeholder{
	"antigravity": {placeholder, "gemini": {placeholder, "codex": {placeholder, "messages": {placeholder, "responses": {placeholder,
	"chat_completions": {placeholder, "images": {placeholder, "videos": {placeholder, "embeddings": {placeholder, "models": {placeholder, "other": {placeholder,
placeholder

var ingressRejectProtocols = map[string]struct{placeholder{
	"google": {placeholder, "anthropic": {placeholder, "openai": {placeholder, "gateway": {placeholder, "other": {placeholder,
placeholder

// ListIngressRejects returns bounded security aggregates, never raw credentials or request bodies.
func (h *OpsHandler) ListIngressRejects(c *gin.Context) {
	if h.opsService == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
placeholder
	page, pageSize := response.ParsePagination(c)
	if pageSize > 200 {
		pageSize = 200
placeholder
	startTime, endTime, err := parseOpsTimeRange(c, "1h")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder
	filter := &service.OpsIngressRejectFilter{Page: page, PageSize: pageSizeplaceholder
	if !startTime.IsZero() {
		filter.StartTime = &startTime
placeholder
	if !endTime.IsZero() {
		filter.EndTime = &endTime
placeholder
	if filter.RejectReason, err = parseIngressRejectEnum(c, "reason", ingressRejectReasons); err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder
	if filter.RouteFamily, err = parseIngressRejectEnum(c, "route_family", ingressRejectRouteFamilies); err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder
	if filter.Protocol, err = parseIngressRejectEnum(c, "protocol", ingressRejectProtocols); err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder
	if raw := strings.TrimSpace(c.Query("client_ip")); raw != "" {
		addr, parseErr := netip.ParseAddr(raw)
		if parseErr != nil {
			response.BadRequest(c, "Invalid client_ip")
			return
	placeholder
		addr = addr.Unmap()
		if addr.Is6() {
			addr = netip.PrefixFrom(addr, 64).Masked().Addr()
	placeholder
		filter.ClientIP = addr.String()
placeholder
	if filter.UserID, err = parseOptionalPositiveID(c, "user_id"); err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder
	if filter.APIKeyID, err = parseOptionalPositiveID(c, "api_key_id"); err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder

	result, err := h.opsService.ListIngressRejects(c.Request.Context(), filter)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, result)
placeholder

func (h *OpsHandler) GetIngressRejectHealth(c *gin.Context) {
	if h.opsService == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
placeholder
	if err := h.opsService.RequireMonitoringEnabled(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, h.opsService.GetIngressRejectHealth())
placeholder

func parseIngressRejectEnum(c *gin.Context, name string, allowed map[string]struct{placeholder) (string, error) {
	value := strings.TrimSpace(c.Query(name))
	if value == "" {
		return "", nil
placeholder
	if _, ok := allowed[value]; !ok {
		return "", &ingressRejectQueryError{message: "Invalid " + nameplaceholder
placeholder
	return value, nil
placeholder

func parseOptionalPositiveID(c *gin.Context, name string) (*int64, error) {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return nil, nil
placeholder
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value <= 0 {
		return nil, &ingressRejectQueryError{message: "Invalid " + nameplaceholder
placeholder
	return &value, nil
placeholder

type ingressRejectQueryError struct{ message string placeholder

func (e *ingressRejectQueryError) Error() string { return e.message placeholder
