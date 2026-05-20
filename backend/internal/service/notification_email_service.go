package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log/slog"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	NotificationEmailEventSubscriptionPurchaseSuccess = "subscription.purchase_success"
	NotificationEmailEventSubscriptionExpiryReminder  = "subscription.expiry_reminder"
	NotificationEmailEventBalanceLow                  = "balance.low"
	NotificationEmailEventBalanceRechargeSuccess      = "balance.recharge_success"

	notificationEmailTemplateKeyPrefix    = "notification_email_template:"
	notificationEmailPreferenceKeyPrefix  = "notification_email_preference:"
	notificationEmailDeliveryKeyPrefix    = "notification_email_delivery:"
	notificationEmailLocaleUserKeyPrefix  = "notification_email_locale:user:"
	notificationEmailLocaleEmailKeyPrefix = "notification_email_locale:email:"
	notificationEmailUnsubscribeSecretKey = "notification_email_unsubscribe_secret"
	notificationEmailDefaultLocale        = "en"
	notificationEmailLocaleChinese        = "zh"
	notificationEmailMaxSubjectLength     = 200
	notificationEmailMaxHTMLLength        = 30000
	notificationEmailUnsubscribeTTL       = 365 * 24 * time.Hour
)

var (
	notificationEmailPlaceholderPattern = regexp.MustCompile(`{{\s*([a-zA-Z][a-zA-Z0-9_]*)\s*placeholderplaceholder`)
	notificationEmailLocales            = []string{notificationEmailDefaultLocale, notificationEmailLocaleChineseplaceholder
	notificationEmailCommonPlaceholders = []string{"site_name", "recipient_name", "recipient_email"placeholder
)

type NotificationEmailService struct {
	settingRepo  SettingRepository
	emailService *EmailService
placeholder

type NotificationEmailEventInfo struct {
	Event        string   `json:"event"`
	Label        string   `json:"label"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	Optional     bool     `json:"optional"`
	Placeholders []string `json:"placeholders"`
placeholder

type NotificationEmailTemplate struct {
	Event        string     `json:"event"`
	Locale       string     `json:"locale"`
	Subject      string     `json:"subject"`
	HTML         string     `json:"html"`
	IsCustom     bool       `json:"is_custom"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
	Placeholders []string   `json:"placeholders"`
placeholder

type NotificationEmailPreview struct {
	Subject string `json:"subject"`
	HTML    string `json:"html"`
placeholder

type NotificationEmailPreviewInput struct {
	Event     string            `json:"event"`
	Locale    string            `json:"locale"`
	Subject   string            `json:"subject"`
	HTML      string            `json:"html"`
	Variables map[string]string `json:"variables,omitempty"`
placeholder

type NotificationEmailSendInput struct {
	Event          string
	Locale         string
	RecipientEmail string
	RecipientName  string
	UserID         int64
	SourceType     string
	SourceID       string
	ReminderKey    string
	Variables      map[string]string
placeholder

type NotificationEmailUnsubscribeResult struct {
	Event string `json:"event"`
	Email string `json:"email"`
	Done  bool   `json:"done"`
placeholder

type notificationEmailStoredTemplate struct {
	Subject   string    `json:"subject"`
	HTML      string    `json:"html"`
	UpdatedAt time.Time `json:"updated_at"`
placeholder

type notificationEmailOfficialTemplate struct {
	Subject string
	HTML    string
placeholder

type notificationEmailUnsubscribeClaims struct {
	Email string `json:"email"`
	Event string `json:"event"`
	Exp   int64  `json:"exp"`
placeholder

func NewNotificationEmailService(settingRepo SettingRepository, emailService *EmailService) *NotificationEmailService {
	return &NotificationEmailService{settingRepo: settingRepo, emailService: emailServiceplaceholder
placeholder

func (s *NotificationEmailService) ListEventInfos() []NotificationEmailEventInfo {
	infos := make([]NotificationEmailEventInfo, 0, len(notificationEmailEventDefinitions))
	for _, event := range notificationEmailEventOrder {
		info := notificationEmailEventDefinitions[event]
		info.Placeholders = append([]string(nil), info.Placeholders...)
		infos = append(infos, info)
placeholder
	return infos
placeholder

func (s *NotificationEmailService) SupportedLocales() []string {
	return append([]string(nil), notificationEmailLocales...)
placeholder

func (s *NotificationEmailService) ListTemplates(ctx context.Context) ([]NotificationEmailTemplate, error) {
	items := make([]NotificationEmailTemplate, 0, len(notificationEmailEventOrder)*len(notificationEmailLocales))
	for _, event := range notificationEmailEventOrder {
		for _, locale := range notificationEmailLocales {
			tmpl, err := s.GetTemplate(ctx, event, locale)
			if err != nil {
				return nil, err
		placeholder
			items = append(items, tmpl)
	placeholder
placeholder
	return items, nil
placeholder

func (s *NotificationEmailService) GetTemplate(ctx context.Context, event, locale string) (NotificationEmailTemplate, error) {
	info, normalizedEvent, err := s.eventInfo(event)
	if err != nil {
		return NotificationEmailTemplate{placeholder, err
placeholder
	normalizedLocale := normalizeNotificationLocale(locale)
	official, ok := notificationEmailOfficialTemplates[normalizedEvent][normalizedLocale]
	if !ok {
		return NotificationEmailTemplate{placeholder, fmt.Errorf("official template not found for %s/%s", normalizedEvent, normalizedLocale)
placeholder

	tmpl := NotificationEmailTemplate{
		Event:        normalizedEvent,
		Locale:       normalizedLocale,
		Subject:      official.Subject,
		HTML:         official.HTML,
		Placeholders: append([]string(nil), info.Placeholders...),
placeholder

	raw, err := s.settingRepo.GetValue(ctx, notificationEmailTemplateKey(normalizedEvent, normalizedLocale))
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return tmpl, nil
	placeholder
		return NotificationEmailTemplate{placeholder, err
placeholder
	if strings.TrimSpace(raw) == "" {
		return tmpl, nil
placeholder

	var stored notificationEmailStoredTemplate
	if err := json.Unmarshal([]byte(raw), &stored); err != nil {
		return NotificationEmailTemplate{placeholder, fmt.Errorf("decode email template override: %w", err)
placeholder
	if err := validateNotificationEmailTemplate(normalizedEvent, stored.Subject, stored.HTML); err != nil {
		return NotificationEmailTemplate{placeholder, err
placeholder
	tmpl.Subject = stored.Subject
	tmpl.HTML = stored.HTML
	tmpl.IsCustom = true
	updatedAt := stored.UpdatedAt
	tmpl.UpdatedAt = &updatedAt
	return tmpl, nil
placeholder

func (s *NotificationEmailService) UpdateTemplate(ctx context.Context, event, locale, subject, htmlBody string) (NotificationEmailTemplate, error) {
	_, normalizedEvent, err := s.eventInfo(event)
	if err != nil {
		return NotificationEmailTemplate{placeholder, err
placeholder
	normalizedLocale := normalizeNotificationLocale(locale)
	if err := validateNotificationEmailTemplate(normalizedEvent, subject, htmlBody); err != nil {
		return NotificationEmailTemplate{placeholder, err
placeholder
	stored := notificationEmailStoredTemplate{
		Subject:   strings.TrimSpace(subject),
		HTML:      htmlBody,
		UpdatedAt: time.Now().UTC(),
placeholder
	payload, err := json.Marshal(stored)
	if err != nil {
		return NotificationEmailTemplate{placeholder, err
placeholder
	if err := s.settingRepo.Set(ctx, notificationEmailTemplateKey(normalizedEvent, normalizedLocale), string(payload)); err != nil {
		return NotificationEmailTemplate{placeholder, err
placeholder
	return s.GetTemplate(ctx, normalizedEvent, normalizedLocale)
placeholder

func (s *NotificationEmailService) RestoreOfficialTemplate(ctx context.Context, event, locale string) (NotificationEmailTemplate, error) {
	_, normalizedEvent, err := s.eventInfo(event)
	if err != nil {
		return NotificationEmailTemplate{placeholder, err
placeholder
	normalizedLocale := normalizeNotificationLocale(locale)
	if err := s.settingRepo.Delete(ctx, notificationEmailTemplateKey(normalizedEvent, normalizedLocale)); err != nil && !errors.Is(err, ErrSettingNotFound) {
		return NotificationEmailTemplate{placeholder, err
placeholder
	return s.GetTemplate(ctx, normalizedEvent, normalizedLocale)
placeholder

func (s *NotificationEmailService) PreviewTemplate(ctx context.Context, input NotificationEmailPreviewInput) (NotificationEmailPreview, error) {
	_, normalizedEvent, err := s.eventInfo(input.Event)
	if err != nil {
		return NotificationEmailPreview{placeholder, err
placeholder
	normalizedLocale := normalizeNotificationLocale(input.Locale)
	subject := input.Subject
	htmlBody := input.HTML
	if strings.TrimSpace(subject) == "" || strings.TrimSpace(htmlBody) == "" {
		tmpl, err := s.GetTemplate(ctx, normalizedEvent, normalizedLocale)
		if err != nil {
			return NotificationEmailPreview{placeholder, err
	placeholder
		if strings.TrimSpace(subject) == "" {
			subject = tmpl.Subject
	placeholder
		if strings.TrimSpace(htmlBody) == "" {
			htmlBody = tmpl.HTML
	placeholder
placeholder
	if err := validateNotificationEmailTemplate(normalizedEvent, subject, htmlBody); err != nil {
		return NotificationEmailPreview{placeholder, err
placeholder
	variables := s.sampleVariables(ctx, normalizedEvent, normalizedLocale)
	for key, value := range input.Variables {
		variables[key] = value
placeholder
	return renderNotificationEmail(normalizedEvent, subject, htmlBody, variables)
placeholder

func (s *NotificationEmailService) Send(ctx context.Context, input NotificationEmailSendInput) error {
	info, normalizedEvent, err := s.eventInfo(input.Event)
	if err != nil {
		return err
placeholder
	recipient := strings.TrimSpace(input.RecipientEmail)
	if recipient == "" {
		return nil
placeholder
	if info.Optional {
		unsubscribed, err := s.IsUnsubscribed(ctx, recipient, normalizedEvent)
		if err != nil {
			return err
	placeholder
		if unsubscribed {
			slog.Info("notification email suppressed by unsubscribe preference", "event", normalizedEvent, "recipient_hash", notificationEmailHash(recipient))
			return nil
	placeholder
placeholder

	locale := normalizeNotificationLocale(input.Locale)
	if strings.TrimSpace(input.Locale) == "" {
		locale = s.ResolveRecipientLocale(ctx, input.UserID, recipient)
placeholder
	tmpl, err := s.GetTemplate(ctx, normalizedEvent, locale)
	if err != nil {
		return err
placeholder
	variables := s.runtimeVariables(ctx, normalizedEvent, locale, input)
	rendered, err := renderNotificationEmail(normalizedEvent, tmpl.Subject, tmpl.HTML, variables)
	if err != nil {
		return err
placeholder

	deliveryKey := notificationEmailDeliveryKey(normalizedEvent, input.SourceType, input.SourceID, recipient, input.ReminderKey)
	if deliveryKey != "" {
		sent, err := s.deliveryExists(ctx, deliveryKey)
		if err != nil {
			return err
	placeholder
		if sent {
			return nil
	placeholder
placeholder

	if s.emailService == nil {
		return errors.New("email service is not configured")
placeholder
	if err := s.emailService.SendEmail(ctx, recipient, rendered.Subject, rendered.HTML); err != nil {
		return err
placeholder
	if deliveryKey != "" {
		_ = s.settingRepo.Set(ctx, deliveryKey, time.Now().UTC().Format(time.RFC3339Nano))
placeholder
	return nil
placeholder

func (s *NotificationEmailService) RememberRecipientLocale(ctx context.Context, userID int64, email, acceptLanguage string) {
	locale := normalizeNotificationLocale(acceptLanguage)
	if strings.TrimSpace(acceptLanguage) == "" || s == nil || s.settingRepo == nil {
		return
placeholder
	if userID > 0 {
		_ = s.settingRepo.Set(ctx, notificationEmailLocaleUserKeyPrefix+strconv.FormatInt(userID, 10), locale)
placeholder
	if emailHash := notificationEmailHash(email); emailHash != "" {
		_ = s.settingRepo.Set(ctx, notificationEmailLocaleEmailKeyPrefix+emailHash, locale)
placeholder
placeholder

func (s *NotificationEmailService) ResolveRecipientLocale(ctx context.Context, userID int64, email string) string {
	if s == nil || s.settingRepo == nil {
		return notificationEmailDefaultLocale
placeholder
	if userID > 0 {
		if locale, err := s.settingRepo.GetValue(ctx, notificationEmailLocaleUserKeyPrefix+strconv.FormatInt(userID, 10)); err == nil && strings.TrimSpace(locale) != "" {
			return normalizeNotificationLocale(locale)
	placeholder
placeholder
	if emailHash := notificationEmailHash(email); emailHash != "" {
		if locale, err := s.settingRepo.GetValue(ctx, notificationEmailLocaleEmailKeyPrefix+emailHash); err == nil && strings.TrimSpace(locale) != "" {
			return normalizeNotificationLocale(locale)
	placeholder
placeholder
	return notificationEmailDefaultLocale
placeholder

func (s *NotificationEmailService) IsUnsubscribed(ctx context.Context, email, event string) (bool, error) {
	info, normalizedEvent, err := s.eventInfo(event)
	if err != nil {
		return false, err
placeholder
	if !info.Optional {
		return false, nil
placeholder
	value, err := s.settingRepo.GetValue(ctx, notificationEmailPreferenceKey(normalizedEvent, email))
	if err != nil {
		if errors.Is(err, ErrSettingNotFound) {
			return false, nil
	placeholder
		return false, err
placeholder
	return strings.EqualFold(strings.TrimSpace(value), "unsubscribed"), nil
placeholder

func (s *NotificationEmailService) Unsubscribe(ctx context.Context, token string) (NotificationEmailUnsubscribeResult, error) {
	claims, err := s.parseUnsubscribeToken(ctx, token)
	if err != nil {
		return NotificationEmailUnsubscribeResult{placeholder, err
placeholder
	info, normalizedEvent, err := s.eventInfo(claims.Event)
	if err != nil {
		return NotificationEmailUnsubscribeResult{placeholder, err
placeholder
	if !info.Optional {
		return NotificationEmailUnsubscribeResult{placeholder, fmt.Errorf("%s is transactional and cannot be unsubscribed", normalizedEvent)
placeholder
	if err := s.settingRepo.Set(ctx, notificationEmailPreferenceKey(normalizedEvent, claims.Email), "unsubscribed"); err != nil {
		return NotificationEmailUnsubscribeResult{placeholder, err
placeholder
	return NotificationEmailUnsubscribeResult{Event: normalizedEvent, Email: claims.Email, Done: trueplaceholder, nil
placeholder

func (s *NotificationEmailService) eventInfo(event string) (NotificationEmailEventInfo, string, error) {
	normalized := strings.ToLower(strings.TrimSpace(event))
	info, ok := notificationEmailEventDefinitions[normalized]
	if !ok {
		return NotificationEmailEventInfo{placeholder, "", fmt.Errorf("unsupported email template event: %s", event)
placeholder
	return info, normalized, nil
placeholder

func (s *NotificationEmailService) sampleVariables(ctx context.Context, event, locale string) map[string]string {
	info := notificationEmailEventDefinitions[event]
	variables := make(map[string]string, len(info.Placeholders))
	for key, value := range notificationEmailSampleVariables(locale) {
		variables[key] = value
placeholder
	variables["site_name"] = s.siteName(ctx)
	if variables["unsubscribe_url"] == "" && info.Optional {
		variables["unsubscribe_url"] = "https://example.com/unsubscribe"
placeholder
	return variables
placeholder

func (s *NotificationEmailService) runtimeVariables(ctx context.Context, event, locale string, input NotificationEmailSendInput) map[string]string {
	variables := s.sampleVariables(ctx, event, locale)
	for key, value := range input.Variables {
		variables[key] = value
placeholder
	variables["site_name"] = s.siteName(ctx)
	variables["recipient_email"] = input.RecipientEmail
	if strings.TrimSpace(input.RecipientName) != "" {
		variables["recipient_name"] = input.RecipientName
placeholder
	if notificationEmailEventDefinitions[event].Optional {
		if unsubscribeURL, err := s.buildUnsubscribeURL(ctx, input.RecipientEmail, event); err == nil {
			variables["unsubscribe_url"] = unsubscribeURL
	placeholder
placeholder
	return variables
placeholder

func (s *NotificationEmailService) siteName(ctx context.Context) string {
	if s == nil || s.settingRepo == nil {
		return defaultSiteName
placeholder
	name, err := s.settingRepo.GetValue(ctx, SettingKeySiteName)
	if err != nil || strings.TrimSpace(name) == "" {
		return defaultSiteName
placeholder
	return strings.TrimSpace(name)
placeholder

func (s *NotificationEmailService) baseURL(ctx context.Context) string {
	if s == nil || s.settingRepo == nil {
		return ""
placeholder
	for _, key := range []string{SettingKeyAPIBaseURL, SettingKeyFrontendURLplaceholder {
		value, err := s.settingRepo.GetValue(ctx, key)
		if err == nil && strings.TrimSpace(value) != "" {
			return strings.TrimRight(strings.TrimSpace(value), "/")
	placeholder
placeholder
	return ""
placeholder

func (s *NotificationEmailService) buildUnsubscribeURL(ctx context.Context, email, event string) (string, error) {
	token, err := s.createUnsubscribeToken(ctx, email, event)
	if err != nil {
		return "", err
placeholder
	path := "/api/v1/settings/email-unsubscribe?token=" + url.QueryEscape(token)
	baseURL := s.baseURL(ctx)
	if baseURL == "" {
		return path, nil
placeholder
	return baseURL + path, nil
placeholder

func (s *NotificationEmailService) createUnsubscribeToken(ctx context.Context, email, event string) (string, error) {
	secret, err := s.unsubscribeSecret(ctx)
	if err != nil {
		return "", err
placeholder
	claims := notificationEmailUnsubscribeClaims{Email: strings.TrimSpace(email), Event: event, Exp: time.Now().Add(notificationEmailUnsubscribeTTL).Unix()placeholder
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
placeholder
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := signNotificationEmailToken(secret, encodedPayload)
	return encodedPayload + "." + signature, nil
placeholder

func (s *NotificationEmailService) parseUnsubscribeToken(ctx context.Context, token string) (notificationEmailUnsubscribeClaims, error) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return notificationEmailUnsubscribeClaims{placeholder, errors.New("invalid unsubscribe token")
placeholder
	secret, err := s.unsubscribeSecret(ctx)
	if err != nil {
		return notificationEmailUnsubscribeClaims{placeholder, err
placeholder
	expected := signNotificationEmailToken(secret, parts[0])
	if !hmac.Equal([]byte(expected), []byte(parts[1])) {
		return notificationEmailUnsubscribeClaims{placeholder, errors.New("invalid unsubscribe token signature")
placeholder
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return notificationEmailUnsubscribeClaims{placeholder, errors.New("invalid unsubscribe token payload")
placeholder
	var claims notificationEmailUnsubscribeClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return notificationEmailUnsubscribeClaims{placeholder, errors.New("invalid unsubscribe token payload")
placeholder
	if strings.TrimSpace(claims.Email) == "" || strings.TrimSpace(claims.Event) == "" {
		return notificationEmailUnsubscribeClaims{placeholder, errors.New("invalid unsubscribe token claims")
placeholder
	if claims.Exp <= time.Now().Unix() {
		return notificationEmailUnsubscribeClaims{placeholder, errors.New("unsubscribe token expired")
placeholder
	return claims, nil
placeholder

func (s *NotificationEmailService) unsubscribeSecret(ctx context.Context) (string, error) {
	secret, err := s.settingRepo.GetValue(ctx, notificationEmailUnsubscribeSecretKey)
	if err == nil && strings.TrimSpace(secret) != "" {
		return strings.TrimSpace(secret), nil
placeholder
	if err != nil && !errors.Is(err, ErrSettingNotFound) {
		return "", err
placeholder
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
placeholder
	secret = base64.RawURLEncoding.EncodeToString(buf)
	if err := s.settingRepo.Set(ctx, notificationEmailUnsubscribeSecretKey, secret); err != nil {
		return "", err
placeholder
	return secret, nil
placeholder

func (s *NotificationEmailService) deliveryExists(ctx context.Context, key string) (bool, error) {
	_, err := s.settingRepo.GetValue(ctx, key)
	if err == nil {
		return true, nil
placeholder
	if errors.Is(err, ErrSettingNotFound) {
		return false, nil
placeholder
	return false, err
placeholder

func validateNotificationEmailTemplate(event, subject, htmlBody string) error {
	subject = strings.TrimSpace(subject)
	if subject == "" {
		return errors.New("email subject cannot be empty")
placeholder
	if len([]rune(subject)) > notificationEmailMaxSubjectLength {
		return fmt.Errorf("email subject cannot exceed %d characters", notificationEmailMaxSubjectLength)
placeholder
	if strings.TrimSpace(htmlBody) == "" {
		return errors.New("email html cannot be empty")
placeholder
	if len([]byte(htmlBody)) > notificationEmailMaxHTMLLength {
		return fmt.Errorf("email html cannot exceed %d bytes", notificationEmailMaxHTMLLength)
placeholder
	allowed := notificationEmailAllowedPlaceholderSet(event)
	for _, placeholder := range notificationEmailPlaceholdersIn(subject + "\n" + htmlBody) {
		if _, ok := allowed[placeholder]; !ok {
			return fmt.Errorf("unsupported placeholder {{%splaceholderplaceholder for event %s", placeholder, event)
	placeholder
placeholder
	return nil
placeholder

func renderNotificationEmail(event, subject, htmlBody string, variables map[string]string) (NotificationEmailPreview, error) {
	if err := validateNotificationEmailTemplate(event, subject, htmlBody); err != nil {
		return NotificationEmailPreview{placeholder, err
placeholder
	renderedSubject, err := renderNotificationEmailString(event, subject, variables, false)
	if err != nil {
		return NotificationEmailPreview{placeholder, err
placeholder
	renderedHTML, err := renderNotificationEmailString(event, htmlBody, variables, true)
	if err != nil {
		return NotificationEmailPreview{placeholder, err
placeholder
	return NotificationEmailPreview{Subject: sanitizeEmailHeader(renderedSubject), HTML: renderedHTMLplaceholder, nil
placeholder

func renderNotificationEmailString(event, raw string, variables map[string]string, escapeHTML bool) (string, error) {
	allowed := notificationEmailAllowedPlaceholderSet(event)
	var renderErr error
	rendered := notificationEmailPlaceholderPattern.ReplaceAllStringFunc(raw, func(match string) string {
		if renderErr != nil {
			return ""
	placeholder
		parts := notificationEmailPlaceholderPattern.FindStringSubmatch(match)
		if len(parts) != 2 {
			return ""
	placeholder
		name := parts[1]
		if _, ok := allowed[name]; !ok {
			renderErr = fmt.Errorf("unsupported placeholder {{%splaceholderplaceholder for event %s", name, event)
			return ""
	placeholder
		value := variables[name]
		if strings.HasSuffix(name, "_url") && !isSafeNotificationEmailURL(value) {
			value = ""
	placeholder
		if escapeHTML {
			return html.EscapeString(value)
	placeholder
		return sanitizeEmailHeader(value)
placeholder)
	if renderErr != nil {
		return "", renderErr
placeholder
	return rendered, nil
placeholder

func notificationEmailAllowedPlaceholderSet(event string) map[string]struct{placeholder {
	info := notificationEmailEventDefinitions[event]
	allowed := make(map[string]struct{placeholder, len(info.Placeholders))
	for _, placeholder := range info.Placeholders {
		allowed[placeholder] = struct{placeholder{placeholder
placeholder
	return allowed
placeholder

func notificationEmailPlaceholdersIn(raw string) []string {
	matches := notificationEmailPlaceholderPattern.FindAllStringSubmatch(raw, -1)
	seen := make(map[string]struct{placeholder, len(matches))
	out := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) != 2 {
			continue
	placeholder
		if _, exists := seen[match[1]]; exists {
			continue
	placeholder
		seen[match[1]] = struct{placeholder{placeholder
		out = append(out, match[1])
placeholder
	return out
placeholder

func normalizeNotificationLocale(raw string) string {
	trimmed := strings.ToLower(strings.TrimSpace(raw))
	if trimmed == "" {
		return notificationEmailDefaultLocale
placeholder
	for _, part := range strings.Split(trimmed, ",") {
		tag := strings.TrimSpace(strings.Split(part, ";")[0])
		if strings.HasPrefix(tag, "zh") || tag == "cn" {
			return notificationEmailLocaleChinese
	placeholder
		if strings.HasPrefix(tag, "en") {
			return notificationEmailDefaultLocale
	placeholder
placeholder
	return notificationEmailDefaultLocale
placeholder

func notificationEmailTemplateKey(event, locale string) string {
	return notificationEmailTemplateKeyPrefix + event + ":" + locale
placeholder

func notificationEmailPreferenceKey(event, email string) string {
	return notificationEmailPreferenceKeyPrefix + event + ":" + notificationEmailHash(email)
placeholder

func notificationEmailDeliveryKey(event, sourceType, sourceID, recipient, reminderKey string) string {
	if strings.TrimSpace(sourceType) == "" || strings.TrimSpace(sourceID) == "" || strings.TrimSpace(recipient) == "" {
		return ""
placeholder
	parts := []string{notificationEmailDeliveryKeyPrefix, event, ":", safeNotificationEmailKeyPart(sourceType), ":", safeNotificationEmailKeyPart(sourceID), ":", notificationEmailHash(recipient)placeholder
	if strings.TrimSpace(reminderKey) != "" {
		parts = append(parts, ":", safeNotificationEmailKeyPart(reminderKey))
placeholder
	return strings.Join(parts, "")
placeholder

func notificationEmailHash(value string) string {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	if trimmed == "" {
		return ""
placeholder
	sum := sha256.Sum256([]byte(trimmed))
	return hex.EncodeToString(sum[:])
placeholder

func safeNotificationEmailKeyPart(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r == '.' {
			builder.WriteRune(r)
	placeholder else {
			builder.WriteRune('_')
	placeholder
placeholder
	return builder.String()
placeholder

func signNotificationEmailToken(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
placeholder

func isSafeNotificationEmailURL(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return true
placeholder
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return false
placeholder
	if parsed.IsAbs() {
		scheme := strings.ToLower(parsed.Scheme)
		return scheme == "http" || scheme == "https" || scheme == "mailto"
placeholder
	return strings.HasPrefix(trimmed, "/")
placeholder

func notificationEmailSampleVariables(locale string) map[string]string {
	if normalizeNotificationLocale(locale) == notificationEmailLocaleChinese {
		return map[string]string{
			"site_name":          defaultSiteName,
			"recipient_name":     "张三",
			"recipient_email":    "user@example.com",
			"subscription_group": "Claude Pro",
			"subscription_days":  "30",
			"expiry_time":        "2026-06-18 12:00",
			"days_remaining":     "3",
			"current_balance":    "12.34",
			"threshold":          "20.00",
			"recharge_url":       "https://example.com/recharge",
			"recharge_amount":    "50.00",
			"order_id":           "1024",
			"unsubscribe_url":    "https://example.com/unsubscribe",
	placeholder
placeholder
	return map[string]string{
		"site_name":          defaultSiteName,
		"recipient_name":     "Alex",
		"recipient_email":    "user@example.com",
		"subscription_group": "Claude Pro",
		"subscription_days":  "30",
		"expiry_time":        "2026-06-18 12:00",
		"days_remaining":     "3",
		"current_balance":    "12.34",
		"threshold":          "20.00",
		"recharge_url":       "https://example.com/recharge",
		"recharge_amount":    "50.00",
		"order_id":           "1024",
		"unsubscribe_url":    "https://example.com/unsubscribe",
placeholder
placeholder

var notificationEmailEventOrder = []string{
	NotificationEmailEventSubscriptionPurchaseSuccess,
	NotificationEmailEventSubscriptionExpiryReminder,
	NotificationEmailEventBalanceLow,
	NotificationEmailEventBalanceRechargeSuccess,
placeholder

var notificationEmailEventDefinitions = map[string]NotificationEmailEventInfo{
	NotificationEmailEventSubscriptionPurchaseSuccess: {
		Event:        NotificationEmailEventSubscriptionPurchaseSuccess,
		Label:        "Subscription purchase success",
		Description:  "Sent after a subscription purchase is fulfilled.",
		Category:     "subscription",
		Optional:     false,
		Placeholders: append(append([]string{placeholder, notificationEmailCommonPlaceholders...), "subscription_group", "subscription_days", "expiry_time", "order_id"),
placeholder,
	NotificationEmailEventSubscriptionExpiryReminder: {
		Event:        NotificationEmailEventSubscriptionExpiryReminder,
		Label:        "Subscription expiry reminder",
		Description:  "Optional reminder sent before an active subscription expires.",
		Category:     "subscription",
		Optional:     true,
		Placeholders: append(append([]string{placeholder, notificationEmailCommonPlaceholders...), "subscription_group", "expiry_time", "days_remaining", "unsubscribe_url"),
placeholder,
	NotificationEmailEventBalanceLow: {
		Event:        NotificationEmailEventBalanceLow,
		Label:        "Low balance alert",
		Description:  "Optional alert sent when balance crosses the configured low-balance threshold.",
		Category:     "billing",
		Optional:     true,
		Placeholders: append(append([]string{placeholder, notificationEmailCommonPlaceholders...), "current_balance", "threshold", "recharge_url", "unsubscribe_url"),
placeholder,
	NotificationEmailEventBalanceRechargeSuccess: {
		Event:        NotificationEmailEventBalanceRechargeSuccess,
		Label:        "Balance recharge success",
		Description:  "Sent after a balance recharge order is fulfilled.",
		Category:     "billing",
		Optional:     false,
		Placeholders: append(append([]string{placeholder, notificationEmailCommonPlaceholders...), "recharge_amount", "current_balance", "order_id"),
placeholder,
placeholder

var notificationEmailOfficialTemplates = map[string]map[string]notificationEmailOfficialTemplate{
	NotificationEmailEventSubscriptionPurchaseSuccess: {
		notificationEmailDefaultLocale: {
			Subject: "[{{site_nameplaceholderplaceholder] Subscription purchase successful",
			HTML: notificationEmailCard("#2563eb", "Subscription activated", `
<p>Hello {{recipient_nameplaceholderplaceholder,</p>
<p>Your subscription for <strong>{{subscription_groupplaceholderplaceholder</strong> has been activated for <strong>{{subscription_daysplaceholderplaceholder</strong> days.</p>
<p>Expiry time: <strong>{{expiry_timeplaceholderplaceholder</strong></p>
<p>Order ID: {{order_idplaceholderplaceholder</p>`),
	placeholder,
		notificationEmailLocaleChinese: {
			Subject: "[{{site_nameplaceholderplaceholder] 订阅购买成功",
			HTML: notificationEmailCard("#2563eb", "订阅已开通", `
<p>{{recipient_nameplaceholderplaceholder，您好：</p>
<p>您的 <strong>{{subscription_groupplaceholderplaceholder</strong> 订阅已成功开通，有效期 <strong>{{subscription_daysplaceholderplaceholder</strong> 天。</p>
<p>到期时间：<strong>{{expiry_timeplaceholderplaceholder</strong></p>
<p>订单号：{{order_idplaceholderplaceholder</p>`),
	placeholder,
placeholder,
	NotificationEmailEventSubscriptionExpiryReminder: {
		notificationEmailDefaultLocale: {
			Subject: "[{{site_nameplaceholderplaceholder] Subscription expires in {{days_remainingplaceholderplaceholder day(s)",
			HTML: notificationEmailCard("#f97316", "Subscription expiry reminder", `
<p>Hello {{recipient_nameplaceholderplaceholder,</p>
<p>Your <strong>{{subscription_groupplaceholderplaceholder</strong> subscription will expire in <strong>{{days_remainingplaceholderplaceholder</strong> day(s).</p>
<p>Expiry time: <strong>{{expiry_timeplaceholderplaceholder</strong></p>
<p class="muted"><a href="{{unsubscribe_urlplaceholderplaceholder">Unsubscribe from optional subscription reminders</a></p>`),
	placeholder,
		notificationEmailLocaleChinese: {
			Subject: "[{{site_nameplaceholderplaceholder] 订阅将在 {{days_remainingplaceholderplaceholder 天后到期",
			HTML: notificationEmailCard("#f97316", "订阅到期提醒", `
<p>{{recipient_nameplaceholderplaceholder，您好：</p>
<p>您的 <strong>{{subscription_groupplaceholderplaceholder</strong> 订阅将在 <strong>{{days_remainingplaceholderplaceholder</strong> 天后到期。</p>
<p>到期时间：<strong>{{expiry_timeplaceholderplaceholder</strong></p>
<p class="muted"><a href="{{unsubscribe_urlplaceholderplaceholder">退订此类订阅提醒</a></p>`),
	placeholder,
placeholder,
	NotificationEmailEventBalanceLow: {
		notificationEmailDefaultLocale: {
			Subject: "[{{site_nameplaceholderplaceholder] Low balance alert",
			HTML: notificationEmailCard("#d97706", "Low balance alert", `
<p>Hello {{recipient_nameplaceholderplaceholder,</p>
<p>Your current balance is <strong>${{current_balanceplaceholderplaceholder</strong>, below the configured alert threshold of <strong>${{thresholdplaceholderplaceholder</strong>.</p>
<p>Please recharge in time to avoid service interruption.</p>
<p><a class="button" href="{{recharge_urlplaceholderplaceholder">Recharge now</a></p>
<p class="muted"><a href="{{unsubscribe_urlplaceholderplaceholder">Unsubscribe from optional balance alerts</a></p>`),
	placeholder,
		notificationEmailLocaleChinese: {
			Subject: "[{{site_nameplaceholderplaceholder] 余额不足提醒",
			HTML: notificationEmailCard("#d97706", "余额不足提醒", `
<p>{{recipient_nameplaceholderplaceholder，您好：</p>
<p>您当前余额为 <strong>${{current_balanceplaceholderplaceholder</strong>，已低于提醒阈值 <strong>${{thresholdplaceholderplaceholder</strong>。</p>
<p>请及时充值以免服务中断。</p>
<p><a class="button" href="{{recharge_urlplaceholderplaceholder">立即充值</a></p>
<p class="muted"><a href="{{unsubscribe_urlplaceholderplaceholder">退订此类余额提醒</a></p>`),
	placeholder,
placeholder,
	NotificationEmailEventBalanceRechargeSuccess: {
		notificationEmailDefaultLocale: {
			Subject: "[{{site_nameplaceholderplaceholder] Balance recharge successful",
			HTML: notificationEmailCard("#16a34a", "Recharge successful", `
<p>Hello {{recipient_nameplaceholderplaceholder,</p>
<p>Your balance recharge of <strong>${{recharge_amountplaceholderplaceholder</strong> has been completed.</p>
<p>Current balance: <strong>${{current_balanceplaceholderplaceholder</strong></p>
<p>Order ID: {{order_idplaceholderplaceholder</p>`),
	placeholder,
		notificationEmailLocaleChinese: {
			Subject: "[{{site_nameplaceholderplaceholder] 余额充值成功",
			HTML: notificationEmailCard("#16a34a", "余额充值成功", `
<p>{{recipient_nameplaceholderplaceholder，您好：</p>
<p>您的余额充值 <strong>${{recharge_amountplaceholderplaceholder</strong> 已完成。</p>
<p>当前余额：<strong>${{current_balanceplaceholderplaceholder</strong></p>
<p>订单号：{{order_idplaceholderplaceholder</p>`),
	placeholder,
placeholder,
placeholder

func notificationEmailCard(accent, title, content string) string {
	return `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 24px; background: #f4f4f5; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; color: #18181b; placeholder
    .container { max-width: 640px; margin: 0 auto; background: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 8px 30px rgba(15, 23, 42, 0.10); placeholder
    .header { background: ` + accent + `; color: #ffffff; padding: 28px 32px; placeholder
    .header h1 { margin: 0; font-size: 24px; line-height: 1.25; placeholder
    .content { padding: 32px; font-size: 15px; line-height: 1.7; placeholder
    .button { display: inline-block; margin-top: 12px; padding: 11px 18px; border-radius: 8px; background: ` + accent + `; color: #ffffff; text-decoration: none; font-weight: 600; placeholder
    .muted { color: #71717a; font-size: 13px; placeholder
    .footer { padding: 18px 32px; background: #fafafa; color: #a1a1aa; font-size: 12px; placeholder
  </style>
</head>
<body>
  <div class="container">
    <div class="header"><h1>` + title + `</h1></div>
    <div class="content">` + content + `</div>
    <div class="footer">This email was sent by {{site_nameplaceholderplaceholder. Please do not reply directly.</div>
  </div>
</body>
</html>`
placeholder
