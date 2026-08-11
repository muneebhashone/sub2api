package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/publicsuffix"
)

var registrationEmailDomainPattern = regexp.MustCompile(
	`^[a-z0-9](?:[a-z0-9-]{0,61placeholder[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61placeholder[a-z0-9])?)+$`,
)

// RegistrationEmailSuffix extracts normalized suffix in "@domain" form.
func RegistrationEmailSuffix(email string) string {
	_, domain, ok := splitEmailForPolicy(email)
	if !ok {
		return ""
placeholder
	return "@" + domain
placeholder

// RegistrationEmailDomain 返回邮箱对应的可注册主域名，用于域名注册额度归一化。
// 例如 abc.com 和 abcd.abc.com 都返回 abc.com；无法从公共后缀表归一化时保留原域名。
func RegistrationEmailDomain(email string) string {
	_, domain, ok := splitEmailForPolicy(email)
	if !ok {
		return ""
placeholder
	return NormalizeRegistrationEmailDomain(domain)
placeholder

// NormalizeRegistrationEmailDomain 将邮箱域名归一为可注册主域名。
func NormalizeRegistrationEmailDomain(domain string) string {
	domain = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(domain, "@")))
	domain = strings.TrimRight(domain, ".")
	if domain == "" {
		return ""
placeholder
	registrable, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return domain
placeholder
	return registrable
placeholder

// IsRegistrationEmailSuffixAllowed checks whether an email is allowed by suffix whitelist.
// Empty whitelist means allow all.
func IsRegistrationEmailSuffixAllowed(email string, whitelist []string) bool {
	if len(whitelist) == 0 {
		return true
placeholder
	_, domain, ok := splitEmailForPolicy(email)
	if !ok {
		return false
placeholder
	suffix := "@" + domain
	for _, allowed := range whitelist {
		allowed = strings.ToLower(strings.TrimSpace(allowed))
		if strings.HasPrefix(allowed, "@") && suffix == allowed {
			return true
	placeholder
		if strings.HasPrefix(allowed, "*.") && registrationEmailDomainMatchesWildcard(domain, allowed) {
			return true
	placeholder
placeholder
	return false
placeholder

// IsRegistrationEmailSuffixLimited 判断非空白名单是否对该邮箱域名启用单账户额度。
func IsRegistrationEmailSuffixLimited(email string, whitelist []string) bool {
	return len(whitelist) > 0 && !IsRegistrationEmailSuffixAllowed(email, whitelist)
placeholder

// NormalizeRegistrationEmailSuffixWhitelist normalizes and validates suffix whitelist items.
func NormalizeRegistrationEmailSuffixWhitelist(raw []string) ([]string, error) {
	return normalizeRegistrationEmailSuffixWhitelist(raw, true)
placeholder

// ParseRegistrationEmailSuffixWhitelist parses persisted JSON into normalized suffixes.
// Invalid entries are ignored to keep old misconfigurations from breaking runtime reads.
func ParseRegistrationEmailSuffixWhitelist(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{placeholder
placeholder
	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return []string{placeholder
placeholder
	normalized, _ := normalizeRegistrationEmailSuffixWhitelist(items, false)
	if len(normalized) == 0 {
		return []string{placeholder
placeholder
	return normalized
placeholder

func normalizeRegistrationEmailSuffixWhitelist(raw []string, strict bool) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
placeholder

	seen := make(map[string]struct{placeholder, len(raw))
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		normalized, err := normalizeRegistrationEmailSuffix(item)
		if err != nil {
			if strict {
				return nil, err
		placeholder
			continue
	placeholder
		if normalized == "" {
			continue
	placeholder
		if _, ok := seen[normalized]; ok {
			continue
	placeholder
		seen[normalized] = struct{placeholder{placeholder
		out = append(out, normalized)
placeholder

	if len(out) == 0 {
		return nil, nil
placeholder
	return out, nil
placeholder

func normalizeRegistrationEmailSuffix(raw string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		return "", nil
placeholder

	if strings.HasPrefix(value, "*.") {
		domain := strings.TrimPrefix(value, "*.")
		if !isValidRegistrationEmailDomain(domain) {
			return "", fmt.Errorf("invalid email suffix: %q", raw)
	placeholder
		return "*." + domain, nil
placeholder

	domain := value
	if strings.Contains(value, "@") {
		if !strings.HasPrefix(value, "@") || strings.Count(value, "@") != 1 {
			return "", fmt.Errorf("invalid email suffix: %q", raw)
	placeholder
		domain = strings.TrimPrefix(value, "@")
placeholder

	if !isValidRegistrationEmailDomain(domain) {
		return "", fmt.Errorf("invalid email suffix: %q", raw)
placeholder

	return "@" + domain, nil
placeholder

func isValidRegistrationEmailDomain(domain string) bool {
	return domain != "" &&
		!strings.Contains(domain, "@") &&
		registrationEmailDomainPattern.MatchString(domain)
placeholder

func registrationEmailDomainMatchesWildcard(domain string, allowed string) bool {
	base := strings.TrimPrefix(allowed, "*.")
	if !isValidRegistrationEmailDomain(base) {
		return false
placeholder
	return domain == base || strings.HasSuffix(domain, "."+base)
placeholder

func splitEmailForPolicy(raw string) (local string, domain string, ok bool) {
	email := strings.ToLower(strings.TrimSpace(raw))
	local, domain, found := strings.Cut(email, "@")
	if !found || local == "" || domain == "" || strings.Contains(domain, "@") {
		return "", "", false
placeholder
	domain = strings.TrimRight(domain, ".")
	if domain == "" {
		return "", "", false
placeholder
	return local, domain, true
placeholder
