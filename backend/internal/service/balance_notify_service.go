package service

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

const (
	emailSendTimeout = 30 * time.Second

	// Quota dimension labels
	quotaDimDaily  = "daily"
	quotaDimWeekly = "weekly"
	quotaDimTotal  = "total"
)

// quotaDimLabels maps dimension names to display labels.
var quotaDimLabels = map[string]string{
	quotaDimDaily:  "日限额 / Daily",
	quotaDimWeekly: "周限额 / Weekly",
	quotaDimTotal:  "总限额 / Total",
placeholder

// AccountQuotaReader provides read access to account quota data.
type AccountQuotaReader interface {
	GetByID(ctx context.Context, id int64) (*Account, error)
placeholder

// BalanceNotifyService handles balance and quota threshold notifications.
type BalanceNotifyService struct {
	emailService *EmailService
	settingRepo  SettingRepository
	accountRepo  AccountQuotaReader
placeholder

// NewBalanceNotifyService creates a new BalanceNotifyService.
func NewBalanceNotifyService(emailService *EmailService, settingRepo SettingRepository, accountRepo AccountQuotaReader) *BalanceNotifyService {
	return &BalanceNotifyService{
		emailService: emailService,
		settingRepo:  settingRepo,
		accountRepo:  accountRepo,
placeholder
placeholder

// CheckBalanceAfterDeduction checks if balance crossed below threshold after deduction.
// oldBalance is the balance before deduction, cost is the amount deducted.
// Notification is sent only on first crossing: oldBalance >= threshold && newBalance < threshold.
func (s *BalanceNotifyService) CheckBalanceAfterDeduction(ctx context.Context, user *User, oldBalance, cost float64) {
	if user == nil || s.emailService == nil || s.settingRepo == nil {
		return
placeholder
	if !user.BalanceNotifyEnabled {
		return
placeholder

	globalEnabled, globalThreshold := s.getBalanceNotifyConfig(ctx)
	if !globalEnabled {
		return
placeholder

	// User custom threshold overrides system default
	threshold := globalThreshold
	if user.BalanceNotifyThreshold != nil {
		threshold = *user.BalanceNotifyThreshold
placeholder
	if threshold <= 0 {
		return
placeholder

	newBalance := oldBalance - cost
	if oldBalance >= threshold && newBalance < threshold {
		siteName := s.getSiteName(ctx)
		recipients := s.collectBalanceNotifyRecipients(user)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("panic in balance notification", "recover", r)
			placeholder
		placeholder()
			s.sendBalanceLowEmails(recipients, user.Username, user.Email, newBalance, threshold, siteName)
	placeholder()
placeholder
placeholder

// quotaDim describes one quota dimension for notification checking.
type quotaDim struct {
	name          string
	enabled       bool
	threshold     float64
	thresholdType string // "fixed" (default) or "percentage"
	oldUsed       float64
	limit         float64
placeholder

// resolvedThreshold returns the effective threshold value.
// For percentage type, it computes threshold = limit * percentage / 100.
func (d quotaDim) resolvedThreshold() float64 {
	if d.thresholdType == "percentage" && d.limit > 0 {
		return d.limit * d.threshold / 100
placeholder
	return d.threshold
placeholder

// buildQuotaDims returns the three quota dimensions for notification checking.
func buildQuotaDims(account *Account) []quotaDim {
	return []quotaDim{
		{quotaDimDaily, account.GetQuotaNotifyDailyEnabled(), account.GetQuotaNotifyDailyThreshold(), account.GetQuotaNotifyDailyThresholdType(), account.GetQuotaDailyUsed(), account.GetQuotaDailyLimit()placeholder,
		{quotaDimWeekly, account.GetQuotaNotifyWeeklyEnabled(), account.GetQuotaNotifyWeeklyThreshold(), account.GetQuotaNotifyWeeklyThresholdType(), account.GetQuotaWeeklyUsed(), account.GetQuotaWeeklyLimit()placeholder,
		{quotaDimTotal, account.GetQuotaNotifyTotalEnabled(), account.GetQuotaNotifyTotalThreshold(), account.GetQuotaNotifyTotalThresholdType(), account.GetQuotaUsed(), account.GetQuotaLimit()placeholder,
placeholder
placeholder

// CheckAccountQuotaAfterIncrement checks if any quota dimension crossed above its notify threshold.
// It fetches real-time quota usage from DB to avoid stale snapshot values.
func (s *BalanceNotifyService) CheckAccountQuotaAfterIncrement(ctx context.Context, account *Account, cost float64) {
	if account == nil || s.emailService == nil || s.settingRepo == nil || cost <= 0 {
		return
placeholder
	if !s.isAccountQuotaNotifyEnabled(ctx) {
		return
placeholder
	adminEmails := s.getAccountQuotaNotifyEmails(ctx)
	if len(adminEmails) == 0 {
		return
placeholder

	freshAccount := s.fetchFreshAccount(ctx, account)
	siteName := s.getSiteName(ctx)
	s.checkQuotaDimCrossings(freshAccount, cost, adminEmails, siteName)
placeholder

// fetchFreshAccount loads the latest account from DB; falls back to the snapshot on error.
func (s *BalanceNotifyService) fetchFreshAccount(ctx context.Context, snapshot *Account) *Account {
	if s.accountRepo == nil {
		return snapshot
placeholder
	fresh, err := s.accountRepo.GetByID(ctx, snapshot.ID)
	if err != nil {
		slog.Warn("failed to fetch fresh account for quota notify, using snapshot",
			"account_id", snapshot.ID, "error", err)
		return snapshot
placeholder
	return fresh
placeholder

// checkQuotaDimCrossings iterates quota dimensions and sends alerts for threshold crossings.
// freshAccount has post-increment values; oldUsed is reconstructed as freshUsed - cost.
func (s *BalanceNotifyService) checkQuotaDimCrossings(freshAccount *Account, cost float64, adminEmails []string, siteName string) {
	for _, dim := range buildQuotaDims(freshAccount) {
		if !dim.enabled || dim.threshold <= 0 {
			continue
	placeholder
		effectiveThreshold := dim.resolvedThreshold()
		if effectiveThreshold <= 0 {
			continue
	placeholder
		// dim.oldUsed is actually the post-increment value from fresh DB data;
		// reconstruct pre-increment value to detect threshold crossing.
		newUsed := dim.oldUsed
		oldUsed := dim.oldUsed - cost
		if oldUsed < effectiveThreshold && newUsed >= effectiveThreshold {
			s.asyncSendQuotaAlert(adminEmails, freshAccount.Name, dim, newUsed, effectiveThreshold, siteName)
	placeholder
placeholder
placeholder

// asyncSendQuotaAlert sends quota alert email in a goroutine with panic recovery.
func (s *BalanceNotifyService) asyncSendQuotaAlert(adminEmails []string, accountName string, dim quotaDim, newUsed, effectiveThreshold float64, siteName string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic in quota notification", "recover", r)
		placeholder
	placeholder()
		s.sendQuotaAlertEmails(adminEmails, accountName, dim.name, newUsed, dim.limit, effectiveThreshold, siteName)
placeholder()
placeholder

// getBalanceNotifyConfig reads global balance notification settings.
func (s *BalanceNotifyService) getBalanceNotifyConfig(ctx context.Context) (enabled bool, threshold float64) {
	keys := []string{SettingKeyBalanceLowNotifyEnabled, SettingKeyBalanceLowNotifyThresholdplaceholder
	settings, err := s.settingRepo.GetMultiple(ctx, keys)
	if err != nil {
		return false, 0
placeholder
	enabled = settings[SettingKeyBalanceLowNotifyEnabled] == "true"
	if v := settings[SettingKeyBalanceLowNotifyThreshold]; v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			threshold = f
	placeholder
placeholder
	return
placeholder

// isAccountQuotaNotifyEnabled checks the global account quota notification toggle.
func (s *BalanceNotifyService) isAccountQuotaNotifyEnabled(ctx context.Context) bool {
	val, err := s.settingRepo.GetValue(ctx, SettingKeyAccountQuotaNotifyEnabled)
	if err != nil {
		return false
placeholder
	return val == "true"
placeholder

// getAccountQuotaNotifyEmails reads admin notification emails from settings,
// filtering out disabled entries. Entries with email="" are resolved to the first admin's email.
func (s *BalanceNotifyService) getAccountQuotaNotifyEmails(ctx context.Context) []string {
	raw, err := s.settingRepo.GetValue(ctx, SettingKeyAccountQuotaNotifyEmails)
	if err != nil || strings.TrimSpace(raw) == "" || raw == "[]" {
		return nil
placeholder

	entries := ParseNotifyEmails(raw)
	if len(entries) == 0 {
		return nil
placeholder

	var recipients []string
	seen := make(map[string]bool)
	for _, entry := range entries {
		if entry.Disabled || !entry.Verified {
			continue
	placeholder
		email := strings.TrimSpace(entry.Email)
		if email == "" {
			continue
	placeholder
		lower := strings.ToLower(email)
		if seen[lower] {
			continue
	placeholder
		seen[lower] = true
		recipients = append(recipients, email)
placeholder
	return recipients
placeholder

// getSiteName reads site name from settings with fallback.
func (s *BalanceNotifyService) getSiteName(ctx context.Context) string {
	name, err := s.settingRepo.GetValue(ctx, SettingKeySiteName)
	if err != nil || name == "" {
		return "Sub2API"
placeholder
	return name
placeholder

// collectBalanceNotifyRecipients returns verified, non-disabled email recipients.
// Only emails with verified=true and disabled=false are included.
func (s *BalanceNotifyService) collectBalanceNotifyRecipients(user *User) []string {
	var recipients []string
	seen := make(map[string]bool)

	for _, entry := range user.BalanceNotifyExtraEmails {
		if entry.Disabled || !entry.Verified {
			continue
	placeholder
		email := strings.TrimSpace(entry.Email)
		if email == "" {
			continue
	placeholder
		lower := strings.ToLower(email)
		if seen[lower] {
			continue
	placeholder
		seen[lower] = true
		recipients = append(recipients, email)
placeholder

	return recipients
placeholder

// sendEmails sends an email to all recipients with shared timeout and error logging.
func (s *BalanceNotifyService) sendEmails(recipients []string, subject, body string, logAttrs ...any) {
	for _, to := range recipients {
		ctx, cancel := context.WithTimeout(context.Background(), emailSendTimeout)
		if err := s.emailService.SendEmail(ctx, to, subject, body); err != nil {
			attrs := append([]any{"to", to, "error", errplaceholder, logAttrs...)
			slog.Error("failed to send notification", attrs...)
	placeholder
		cancel()
placeholder
placeholder

// sendBalanceLowEmails sends balance low notification to all recipients.
func (s *BalanceNotifyService) sendBalanceLowEmails(recipients []string, userName, userEmail string, balance, threshold float64, siteName string) {
	displayName := userName
	if displayName == "" {
		displayName = userEmail
placeholder
	subject := fmt.Sprintf("[%s] 余额不足提醒 / Balance Low Alert", siteName)
	body := s.buildBalanceLowEmailBody(html.EscapeString(displayName), balance, threshold, html.EscapeString(siteName))
	s.sendEmails(recipients, subject, body, "user_email", userEmail, "balance", balance)
placeholder

// sendQuotaAlertEmails sends quota alert notification to admin emails.
func (s *BalanceNotifyService) sendQuotaAlertEmails(adminEmails []string, accountName, dimension string, used, limit, threshold float64, siteName string) {
	dimLabel := quotaDimLabels[dimension]
	if dimLabel == "" {
		dimLabel = dimension
placeholder

	subject := fmt.Sprintf("[%s] 账号限额告警 / Account Quota Alert - %s", siteName, accountName)
	body := s.buildQuotaAlertEmailBody(html.EscapeString(accountName), html.EscapeString(dimLabel), used, limit, threshold, html.EscapeString(siteName))
	s.sendEmails(adminEmails, subject, body, "account", accountName, "dimension", dimension)
placeholder

// buildBalanceLowEmailBody builds HTML email for balance low notification.
// Lines exceed 30 due to inline HTML template (not splittable).
func (s *BalanceNotifyService) buildBalanceLowEmailBody(userName string, balance, threshold float64, siteName string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f5f5f5; margin: 0; padding: 20px; placeholder
        .container { max-width: 600px; margin: 0 auto; background-color: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); placeholder
        .header { background: linear-gradient(135deg, #f59e0b 0%%, #d97706 100%%); color: white; padding: 30px; text-align: center; placeholder
        .header h1 { margin: 0; font-size: 24px; placeholder
        .content { padding: 40px 30px; text-align: center; placeholder
        .balance { font-size: 36px; font-weight: bold; color: #dc2626; margin: 20px 0; placeholder
        .info { color: #666; font-size: 14px; line-height: 1.6; margin-top: 20px; placeholder
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; color: #999; font-size: 12px; placeholder
    </style>
</head>
<body>
    <div class="container">
        <div class="header"><h1>%s</h1></div>
        <div class="content">
            <p style="font-size: 18px; color: #333;">%s，您的余额不足</p>
            <p style="color: #666;">Dear %s, your balance is running low</p>
            <div class="balance">$%.2f</div>
            <div class="info">
                <p>您的账户余额已低于提醒阈值 <strong>$%.2f</strong>。</p>
                <p>Your account balance has fallen below the alert threshold of <strong>$%.2f</strong>.</p>
                <p>请及时充值以免服务中断。</p>
                <p>Please top up to avoid service interruption.</p>
            </div>
        </div>
        <div class="footer"><p>此邮件由系统自动发送，请勿回复。</p></div>
    </div>
</body>
</html>`, siteName, userName, userName, balance, threshold, threshold)
placeholder

// buildQuotaAlertEmailBody builds HTML email for account quota alert.
// Lines exceed 30 due to inline HTML template (not splittable).
func (s *BalanceNotifyService) buildQuotaAlertEmailBody(accountName, dimLabel string, used, limit, threshold float64, siteName string) string {
	limitStr := fmt.Sprintf("$%.2f", limit)
	if limit <= 0 {
		limitStr = "无限制 / Unlimited"
placeholder
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f5f5f5; margin: 0; padding: 20px; placeholder
        .container { max-width: 600px; margin: 0 auto; background-color: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); placeholder
        .header { background: linear-gradient(135deg, #ef4444 0%%, #dc2626 100%%); color: white; padding: 30px; text-align: center; placeholder
        .header h1 { margin: 0; font-size: 24px; placeholder
        .content { padding: 40px 30px; placeholder
        .metric { display: flex; justify-content: space-between; padding: 12px 0; border-bottom: 1px solid #eee; placeholder
        .metric-label { color: #666; placeholder
        .metric-value { font-weight: bold; color: #333; placeholder
        .info { color: #666; font-size: 14px; line-height: 1.6; margin-top: 20px; text-align: center; placeholder
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; color: #999; font-size: 12px; placeholder
    </style>
</head>
<body>
    <div class="container">
        <div class="header"><h1>%s</h1></div>
        <div class="content">
            <p style="font-size: 18px; color: #333; text-align: center;">账号限额告警 / Account Quota Alert</p>
            <div class="metric"><span class="metric-label">账号 / Account</span><span class="metric-value">%s</span></div>
            <div class="metric"><span class="metric-label">维度 / Dimension</span><span class="metric-value">%s</span></div>
            <div class="metric"><span class="metric-label">已使用 / Used</span><span class="metric-value">$%.2f</span></div>
            <div class="metric"><span class="metric-label">限额 / Limit</span><span class="metric-value">%s</span></div>
            <div class="metric"><span class="metric-label">告警阈值 / Threshold</span><span class="metric-value">$%.2f</span></div>
            <div class="info">
                <p>账号配额用量已达到告警阈值，请及时关注。</p>
                <p>Account quota usage has reached the alert threshold.</p>
            </div>
        </div>
        <div class="footer"><p>此邮件由系统自动发送，请勿回复。</p></div>
    </div>
</body>
</html>`, siteName, accountName, dimLabel, used, limitStr, threshold)
placeholder

// parseJSONStringArray parses a JSON string array, returns nil on error.
func parseJSONStringArray(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "[]" {
		return nil
placeholder
	var result []string
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil
placeholder
	return result
placeholder
