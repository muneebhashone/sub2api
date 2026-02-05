package admin

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// ExportData exports proxy-only data for migration.
func (h *ProxyHandler) ExportData(c *gin.Context) {
	ctx := c.Request.Context()

	selectedIDs, err := parseProxyIDs(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder

	var proxies []service.Proxy
	if len(selectedIDs) > 0 {
		proxies, err = h.getProxiesByIDs(ctx, selectedIDs)
		if err != nil {
			response.ErrorFrom(c, err)
			return
	placeholder
placeholder else {
		protocol := c.Query("protocol")
		status := c.Query("status")
		search := strings.TrimSpace(c.Query("search"))
		if len(search) > 100 {
			search = search[:100]
	placeholder

		proxies, err = h.listProxiesFiltered(ctx, protocol, status, search)
		if err != nil {
			response.ErrorFrom(c, err)
			return
	placeholder
placeholder

	dataProxies := make([]DataProxy, 0, len(proxies))
	for i := range proxies {
		p := proxies[i]
		key := buildProxyKey(p.Protocol, p.Host, p.Port, p.Username, p.Password)
		dataProxies = append(dataProxies, DataProxy{
			ProxyKey: key,
			Name:     p.Name,
			Protocol: p.Protocol,
			Host:     p.Host,
			Port:     p.Port,
			Username: p.Username,
			Password: p.Password,
			Status:   p.Status,
	placeholder)
placeholder

	payload := DataPayload{
		Type:       dataType,
		Version:    dataVersion,
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		Proxies:    dataProxies,
		Accounts:   []DataAccount{placeholder,
placeholder

	response.Success(c, payload)
placeholder

// ImportData imports proxy-only data for migration.
func (h *ProxyHandler) ImportData(c *gin.Context) {
	type ProxyImportRequest struct {
		Data DataPayload `json:"data"`
placeholder

	var req ProxyImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	if err := validateDataHeader(req.Data); err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder

	ctx := c.Request.Context()
	result := DataImportResult{placeholder

	existingProxies, err := h.listProxiesFiltered(ctx, "", "", "")
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	proxyByKey := make(map[string]service.Proxy, len(existingProxies))
	for i := range existingProxies {
		p := existingProxies[i]
		key := buildProxyKey(p.Protocol, p.Host, p.Port, p.Username, p.Password)
		proxyByKey[key] = p
placeholder

	latencyProbeIDs := make([]int64, 0, len(req.Data.Proxies))
	for i := range req.Data.Proxies {
		item := req.Data.Proxies[i]
		key := item.ProxyKey
		if key == "" {
			key = buildProxyKey(item.Protocol, item.Host, item.Port, item.Username, item.Password)
	placeholder

		if err := validateDataProxy(item); err != nil {
			result.ProxyFailed++
			result.Errors = append(result.Errors, DataImportError{
				Kind:     "proxy",
				Name:     item.Name,
				ProxyKey: key,
				Message:  err.Error(),
		placeholder)
			continue
	placeholder

		if existing, ok := proxyByKey[key]; ok {
			result.ProxyReused++
			if item.Status != "" && item.Status != existing.Status {
				if _, err := h.adminService.UpdateProxy(ctx, existing.ID, &service.UpdateProxyInput{Status: item.Statusplaceholder); err != nil {
					result.Errors = append(result.Errors, DataImportError{
						Kind:     "proxy",
						Name:     item.Name,
						ProxyKey: key,
						Message:  "update status failed: " + err.Error(),
				placeholder)
			placeholder
		placeholder
			latencyProbeIDs = append(latencyProbeIDs, existing.ID)
			continue
	placeholder

		created, err := h.adminService.CreateProxy(ctx, &service.CreateProxyInput{
			Name:     defaultProxyName(item.Name),
			Protocol: item.Protocol,
			Host:     item.Host,
			Port:     item.Port,
			Username: item.Username,
			Password: item.Password,
	placeholder)
		if err != nil {
			result.ProxyFailed++
			result.Errors = append(result.Errors, DataImportError{
				Kind:     "proxy",
				Name:     item.Name,
				ProxyKey: key,
				Message:  err.Error(),
		placeholder)
			continue
	placeholder
		result.ProxyCreated++
		proxyByKey[key] = *created

		if item.Status != "" && item.Status != created.Status {
			if _, err := h.adminService.UpdateProxy(ctx, created.ID, &service.UpdateProxyInput{Status: item.Statusplaceholder); err != nil {
				result.Errors = append(result.Errors, DataImportError{
					Kind:     "proxy",
					Name:     item.Name,
					ProxyKey: key,
					Message:  "update status failed: " + err.Error(),
			placeholder)
		placeholder
	placeholder
		latencyProbeIDs = append(latencyProbeIDs, created.ID)
placeholder

	if len(latencyProbeIDs) > 0 {
		ids := append([]int64(nil), latencyProbeIDs...)
		go func() {
			for _, id := range ids {
				_, _ = h.adminService.TestProxy(context.Background(), id)
		placeholder
	placeholder()
placeholder

	response.Success(c, result)
placeholder

func (h *ProxyHandler) getProxiesByIDs(ctx context.Context, ids []int64) ([]service.Proxy, error) {
	out := make([]service.Proxy, 0, len(ids))
	for _, id := range ids {
		proxy, err := h.adminService.GetProxy(ctx, id)
		if err != nil {
			return nil, err
	placeholder
		if proxy == nil {
			continue
	placeholder
		out = append(out, *proxy)
placeholder
	return out, nil
placeholder

func parseProxyIDs(c *gin.Context) ([]int64, error) {
	values := c.QueryArray("ids")
	if len(values) == 0 {
		raw := strings.TrimSpace(c.Query("ids"))
		if raw != "" {
			values = []string{rawplaceholder
	placeholder
placeholder
	if len(values) == 0 {
		return nil, nil
placeholder

	ids := make([]int64, 0, len(values))
	for _, item := range values {
		for _, part := range strings.Split(item, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
		placeholder
			id, err := strconv.ParseInt(part, 10, 64)
			if err != nil || id <= 0 {
				return nil, fmt.Errorf("invalid proxy id: %s", part)
		placeholder
			ids = append(ids, id)
	placeholder
placeholder
	return ids, nil
placeholder

func (h *ProxyHandler) listProxiesFiltered(ctx context.Context, protocol, status, search string) ([]service.Proxy, error) {
	page := 1
	pageSize := dataPageCap
	var out []service.Proxy
	for {
		items, total, err := h.adminService.ListProxies(ctx, page, pageSize, protocol, status, search)
		if err != nil {
			return nil, err
	placeholder
		out = append(out, items...)
		if len(out) >= int(total) || len(items) == 0 {
			break
	placeholder
		page++
placeholder
	return out, nil
placeholder
