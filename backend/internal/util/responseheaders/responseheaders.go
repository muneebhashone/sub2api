package responseheaders

import (
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

// defaultAllowed 定义允许透传的响应头白名单
// 注意：以下头部由 Go HTTP 包自动处理，不应手动设置：
//   - content-length: 由 ResponseWriter 根据实际写入数据自动设置
//   - transfer-encoding: 由 HTTP 库根据需要自动添加/移除
//   - connection: 由 HTTP 库管理连接复用
var defaultAllowed = map[string]struct{placeholder{
	"content-type":                   {placeholder,
	"content-encoding":               {placeholder,
	"content-language":               {placeholder,
	"cache-control":                  {placeholder,
	"etag":                           {placeholder,
	"last-modified":                  {placeholder,
	"expires":                        {placeholder,
	"vary":                           {placeholder,
	"date":                           {placeholder,
	"x-request-id":                   {placeholder,
	"x-ratelimit-limit-requests":     {placeholder,
	"x-ratelimit-limit-tokens":       {placeholder,
	"x-ratelimit-remaining-requests": {placeholder,
	"x-ratelimit-remaining-tokens":   {placeholder,
	"x-ratelimit-reset-requests":     {placeholder,
	"x-ratelimit-reset-tokens":       {placeholder,
	"retry-after":                    {placeholder,
	"location":                       {placeholder,
placeholder

// hopByHopHeaders 是跳过的 hop-by-hop 头部，这些头部由 HTTP 库自动处理
var hopByHopHeaders = map[string]struct{placeholder{
	"content-length":    {placeholder,
	"transfer-encoding": {placeholder,
	"connection":        {placeholder,
placeholder

func FilterHeaders(src http.Header, cfg config.ResponseHeaderConfig) http.Header {
	allowed := make(map[string]struct{placeholder, len(defaultAllowed)+len(cfg.AdditionalAllowed))
	for key := range defaultAllowed {
		allowed[key] = struct{placeholder{placeholder
placeholder
	for _, key := range cfg.AdditionalAllowed {
		normalized := strings.ToLower(strings.TrimSpace(key))
		if normalized == "" {
			continue
	placeholder
		allowed[normalized] = struct{placeholder{placeholder
placeholder

	forceRemove := make(map[string]struct{placeholder, len(cfg.ForceRemove))
	for _, key := range cfg.ForceRemove {
		normalized := strings.ToLower(strings.TrimSpace(key))
		if normalized == "" {
			continue
	placeholder
		forceRemove[normalized] = struct{placeholder{placeholder
placeholder

	filtered := make(http.Header, len(src))
	for key, values := range src {
		lower := strings.ToLower(key)
		if _, blocked := forceRemove[lower]; blocked {
			continue
	placeholder
		if _, ok := allowed[lower]; !ok {
			continue
	placeholder
		// 跳过 hop-by-hop 头部，这些由 HTTP 库自动处理
		if _, isHopByHop := hopByHopHeaders[lower]; isHopByHop {
			continue
	placeholder
		for _, value := range values {
			filtered.Add(key, value)
	placeholder
placeholder
	return filtered
placeholder

func WriteFilteredHeaders(dst http.Header, src http.Header, cfg config.ResponseHeaderConfig) {
	filtered := FilterHeaders(src, cfg)
	for key, values := range filtered {
		for _, value := range values {
			dst.Add(key, value)
	placeholder
placeholder
placeholder
