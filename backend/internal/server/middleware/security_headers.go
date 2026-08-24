package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
)

const (
	// CSPNonceKey is the context key for storing the CSP nonce
	CSPNonceKey = "csp_nonce"
	// NonceTemplate is the placeholder in CSP policy for nonce
	NonceTemplate = "__CSP_NONCE__"
	// CloudflareInsightsDomain is the domain for Cloudflare Web Analytics
	CloudflareInsightsDomain = "https://static.cloudflareinsights.com"
	// TencentCaptchaDomain is the Tencent Captcha 2.0 Web SDK domain (Chinese mainland site).
	TencentCaptchaDomain = "https://turing.captcha.qcloud.com"
	// TencentCaptchaStaticDomain is the Tencent Captcha static asset domain.
	TencentCaptchaStaticDomain = "https://*.captcha.gtimg.com"
	// TencentCaptchaCDNDomain 是天御国内站的核心 JS CDN 主机：
	// 入口脚本 TJCaptcha.js 会再从这里加载 /1/tgJCap.*.js，缺失时会被 script-src 拦截。
	TencentCaptchaCDNDomain = "https://turing.captcha.gtimg.com"
	// TencentCaptchaGlobalDomain 是天御国际站的 Web SDK 与验证弹窗 iframe 主机。
	TencentCaptchaGlobalDomain = "https://ca.turing.captcha.qcloud.com"
	// TencentCaptchaGlobalCDNDomain 是天御国际站的核心 JS CDN 主机。
	TencentCaptchaGlobalCDNDomain = "https://global.turing.captcha.gtimg.com"
	// TencentCaptchaPrehandleDomain 是天御 SDK 动态预处理脚本与预处理接口主机。
	TencentCaptchaPrehandleDomain = "https://www.tycaptcha.com"
	// TencentCaptchaJQueryDomain 是国内站入口脚本动态加载的 jQuery CDN 主机。
	TencentCaptchaJQueryDomain = "https://cloudcache.tencentcs.com"
	// TencentCaptchaRceDomain 是国际站风控校验接口主机。
	TencentCaptchaRceDomain = "https://rce.tencentrio.com"
	// TencentCaptchaWorkerSource 是天御国际站创建验证码 Web Worker 时使用的来源。
	TencentCaptchaWorkerSource = "blob:"
	// StripeDomain is the domain for Stripe.js SDK
	StripeDomain = "https://*.stripe.com"
	// AirwallexStaticDomain 是 Airwallex 生产环境 SDK 脚本域名。
	AirwallexStaticDomain = "https://static.airwallex.com"
	// AirwallexCheckoutDomain 是 Airwallex 生产环境收银台元素和 iframe 域名。
	AirwallexCheckoutDomain = "https://checkout.airwallex.com"
	// AirwallexDemoStaticDomain 是 Airwallex 沙箱环境 SDK 脚本域名。
	AirwallexDemoStaticDomain = "https://static-demo.airwallex.com"
	// AirwallexDemoCheckoutDomain 是 Airwallex 沙箱环境收银台元素和 iframe 域名。
	AirwallexDemoCheckoutDomain = "https://checkout-demo.airwallex.com"
)

var requiredCSPDirectiveValues = []struct {
	directive string
	value     string
placeholder{
	// 插件配置 UI 使用同源 iframe；目标响应仍必须显式放开 X-Frame-Options，
	// 因此这里只允许 'self' 不会使其他默认 DENY 的管理/API 页面可被嵌入。
	{"frame-src", "'self'"placeholder,
	{"script-src", CloudflareInsightsDomainplaceholder,
	{"script-src", TencentCaptchaDomainplaceholder,
	{"frame-src", TencentCaptchaDomainplaceholder,
	{"style-src", TencentCaptchaStaticDomainplaceholder,
	{"script-src", TencentCaptchaCDNDomainplaceholder,
	{"script-src", TencentCaptchaGlobalDomainplaceholder,
	{"script-src", TencentCaptchaGlobalCDNDomainplaceholder,
	{"script-src", TencentCaptchaPrehandleDomainplaceholder,
	{"script-src", TencentCaptchaJQueryDomainplaceholder,
	{"connect-src", TencentCaptchaDomainplaceholder,
	{"connect-src", TencentCaptchaPrehandleDomainplaceholder,
	{"connect-src", TencentCaptchaRceDomainplaceholder,
	{"frame-src", TencentCaptchaGlobalDomainplaceholder,
	{"frame-src", TencentCaptchaPrehandleDomainplaceholder,
	{"worker-src", TencentCaptchaWorkerSourceplaceholder,
	{"script-src", StripeDomainplaceholder,
	{"frame-src", StripeDomainplaceholder,
	{"script-src", AirwallexStaticDomainplaceholder,
	{"script-src", AirwallexCheckoutDomainplaceholder,
	{"style-src", AirwallexStaticDomainplaceholder,
	{"style-src", AirwallexCheckoutDomainplaceholder,
	{"frame-src", AirwallexCheckoutDomainplaceholder,
	{"script-src", AirwallexDemoStaticDomainplaceholder,
	{"script-src", AirwallexDemoCheckoutDomainplaceholder,
	{"style-src", AirwallexDemoStaticDomainplaceholder,
	{"style-src", AirwallexDemoCheckoutDomainplaceholder,
	{"frame-src", AirwallexDemoCheckoutDomainplaceholder,
placeholder

// GenerateNonce generates a cryptographically secure random nonce.
// 返回 error 以确保调用方在 crypto/rand 失败时能正确降级。
func GenerateNonce() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate CSP nonce: %w", err)
placeholder
	return base64.StdEncoding.EncodeToString(b), nil
placeholder

// GetNonceFromContext retrieves the CSP nonce from gin context
func GetNonceFromContext(c *gin.Context) string {
	if nonce, exists := c.Get(CSPNonceKey); exists {
		if s, ok := nonce.(string); ok {
			return s
	placeholder
placeholder
	return ""
placeholder

// SecurityHeaders sets baseline security headers for all responses.
// getFrameSrcOrigins is an optional function that returns extra origins to inject into frame-src;
// pass nil to disable dynamic frame-src injection.
func SecurityHeaders(cfg config.CSPConfig, getFrameSrcOrigins func() []string) gin.HandlerFunc {
	policy := strings.TrimSpace(cfg.Policy)
	if policy == "" {
		policy = config.DefaultCSPPolicy
placeholder

	// Enhance policy with required directives (nonce placeholder and Cloudflare Insights)
	policy = enhanceCSPPolicy(policy)

	return func(c *gin.Context) {
		finalPolicy := policy
		if getFrameSrcOrigins != nil {
			for _, origin := range getFrameSrcOrigins() {
				if origin != "" {
					finalPolicy = addToDirective(finalPolicy, "frame-src", origin)
			placeholder
		placeholder
	placeholder

		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		if isAPIRoutePath(c) {
			c.Next()
			return
	placeholder

		if cfg.Enabled {
			// Generate nonce for this request
			nonce, err := GenerateNonce()
			if err != nil {
				// crypto/rand 失败时降级为无 nonce 的 CSP 策略
				log.Printf("[SecurityHeaders] %v — 降级为无 nonce 的 CSP", err)
				c.Header("Content-Security-Policy", strings.ReplaceAll(finalPolicy, NonceTemplate, "'unsafe-inline'"))
		placeholder else {
				c.Set(CSPNonceKey, nonce)
				c.Header("Content-Security-Policy", strings.ReplaceAll(finalPolicy, NonceTemplate, "'nonce-"+nonce+"'"))
		placeholder
	placeholder
		c.Next()
placeholder
placeholder

func isAPIRoutePath(c *gin.Context) bool {
	if c == nil || c.Request == nil || c.Request.URL == nil {
		return false
placeholder
	path := c.Request.URL.Path
	return strings.HasPrefix(path, "/v1/") ||
		strings.HasPrefix(path, "/v1beta/") ||
		strings.HasPrefix(path, "/antigravity/") ||
		strings.HasPrefix(path, "/responses") ||
		strings.HasPrefix(path, "/images")
placeholder

// enhanceCSPPolicy 确保 CSP 策略包含 nonce 支持和运行时组件必需域名。
// 这样旧配置文件没有及时补域名时，验证码和支付组件仍能正常加载。
func enhanceCSPPolicy(policy string) string {
	// Add nonce placeholder to script-src if not present
	if !strings.Contains(policy, NonceTemplate) && !strings.Contains(policy, "'nonce-") {
		policy = addToDirective(policy, "script-src", NonceTemplate)
placeholder

	for _, required := range requiredCSPDirectiveValues {
		if !directiveHasValue(policy, required.directive, required.value) {
			policy = addToDirective(policy, required.directive, required.value)
	placeholder
placeholder

	return policy
placeholder

func directiveHasValue(policy, directive, value string) bool {
	for _, rawDirective := range strings.Split(policy, ";") {
		fields := strings.Fields(strings.TrimSpace(rawDirective))
		if len(fields) == 0 || fields[0] != directive {
			continue
	placeholder
		for _, field := range fields[1:] {
			if field == value {
				return true
		placeholder
	placeholder
		return false
placeholder
	return false
placeholder

// addToDirective adds a value to a specific CSP directive.
// If the directive doesn't exist, it will be added after default-src.
func addToDirective(policy, directive, value string) string {
	if end, ok := cspDirectiveEnd(policy, directive); ok {
		return policy[:end] + " " + value + policy[end:]
placeholder
	trimmed := strings.TrimSpace(policy)
	if trimmed == "" {
		return newCSPDirective(directive, value)
placeholder
	if !strings.HasSuffix(trimmed, ";") {
		trimmed += ";"
placeholder
	return trimmed + " " + newCSPDirective(directive, value)
placeholder

func cspDirectiveEnd(policy, directive string) (int, bool) {
	start := 0
	for start <= len(policy) {
		end := len(policy)
		if relativeEnd := strings.IndexByte(policy[start:], ';'); relativeEnd >= 0 {
			end = start + relativeEnd
	placeholder
		fields := strings.Fields(policy[start:end])
		if len(fields) > 0 && fields[0] == directive {
			return end, true
	placeholder
		if end == len(policy) {
			break
	placeholder
		start = end + 1
placeholder
	return 0, false
placeholder

func newCSPDirective(directive, value string) string {
	if value == "'self'" {
		return directive + " 'self';"
placeholder
	return directive + " 'self' " + value + ";"
placeholder
