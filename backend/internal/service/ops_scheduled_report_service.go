package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

const (
	opsScheduledReportJobName = "ops_scheduled_reports"

	opsScheduledReportLeaderLockKeyDefault = "ops:scheduled_reports:leader"
	opsScheduledReportLeaderLockTTLDefault = 5 * time.Minute

	opsScheduledReportLastRunKeyPrefix = "ops:scheduled_reports:last_run:"

	opsScheduledReportTickInterval = 1 * time.Minute
)

var opsScheduledReportCronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

var opsScheduledReportReleaseScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
end
return 0
`)

type OpsScheduledReportService struct {
	opsService   *OpsService
	userService  *UserService
	emailService *EmailService
	redisClient  *redis.Client
	cfg          *config.Config

	instanceID string
	loc        *time.Location

	distributedLockOn bool
	warnNoRedisOnce   sync.Once

	startOnce sync.Once
	stopOnce  sync.Once
	stopCtx   context.Context
	stop      context.CancelFunc
	wg        sync.WaitGroup
placeholder

func NewOpsScheduledReportService(
	opsService *OpsService,
	userService *UserService,
	emailService *EmailService,
	redisClient *redis.Client,
	cfg *config.Config,
) *OpsScheduledReportService {
	lockOn := cfg == nil || strings.TrimSpace(cfg.RunMode) != config.RunModeSimple

	loc := time.Local
	if cfg != nil && strings.TrimSpace(cfg.Timezone) != "" {
		if parsed, err := time.LoadLocation(strings.TrimSpace(cfg.Timezone)); err == nil && parsed != nil {
			loc = parsed
	placeholder
placeholder
	return &OpsScheduledReportService{
		opsService:   opsService,
		userService:  userService,
		emailService: emailService,
		redisClient:  redisClient,
		cfg:          cfg,

		instanceID:        uuid.NewString(),
		loc:               loc,
		distributedLockOn: lockOn,
		warnNoRedisOnce:   sync.Once{placeholder,
		startOnce:         sync.Once{placeholder,
		stopOnce:          sync.Once{placeholder,
		stopCtx:           nil,
		stop:              nil,
		wg:                sync.WaitGroup{placeholder,
placeholder
placeholder

func (s *OpsScheduledReportService) Start() {
	s.StartWithContext(context.Background())
placeholder

func (s *OpsScheduledReportService) StartWithContext(ctx context.Context) {
	if s == nil {
		return
placeholder
	if ctx == nil {
		ctx = context.Background()
placeholder
	if s.cfg != nil && !s.cfg.Ops.Enabled {
		return
placeholder
	if s.opsService == nil || s.emailService == nil {
		return
placeholder

	s.startOnce.Do(func() {
		s.stopCtx, s.stop = context.WithCancel(ctx)
		s.wg.Add(1)
		go s.run()
placeholder)
placeholder

func (s *OpsScheduledReportService) Stop() {
	if s == nil {
		return
placeholder
	s.stopOnce.Do(func() {
		if s.stop != nil {
			s.stop()
	placeholder
placeholder)
	s.wg.Wait()
placeholder

func (s *OpsScheduledReportService) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(opsScheduledReportTickInterval)
	defer ticker.Stop()

	s.runOnce()
	for {
		select {
		case <-ticker.C:
			s.runOnce()
		case <-s.stopCtx.Done():
			return
	placeholder
placeholder
placeholder

func (s *OpsScheduledReportService) runOnce() {
	if s == nil || s.opsService == nil || s.emailService == nil {
		return
placeholder

	startedAt := time.Now().UTC()
	runAt := startedAt

	ctx, cancel := context.WithTimeout(s.stopCtx, 60*time.Second)
	defer cancel()

	// Respect ops monitoring enabled switch.
	if !s.opsService.IsMonitoringEnabled(ctx) {
		return
placeholder

	release, ok := s.tryAcquireLeaderLock(ctx)
	if !ok {
		return
placeholder
	if release != nil {
		defer release()
placeholder

	now := time.Now()
	if s.loc != nil {
		now = now.In(s.loc)
placeholder

	reports := s.listScheduledReports(ctx, now)
	if len(reports) == 0 {
		return
placeholder

	reportsTotal := len(reports)
	reportsDue := 0
	sentAttempts := 0

	for _, report := range reports {
		if report == nil || !report.Enabled {
			continue
	placeholder
		if report.NextRunAt.After(now) {
			continue
	placeholder
		reportsDue++

		attempts, err := s.runReport(ctx, report, now)
		if err != nil {
			s.recordHeartbeatError(runAt, time.Since(startedAt), err)
			return
	placeholder
		sentAttempts += attempts
placeholder

	result := truncateString(fmt.Sprintf("reports=%d due=%d send_attempts=%d", reportsTotal, reportsDue, sentAttempts), 2048)
	s.recordHeartbeatSuccess(runAt, time.Since(startedAt), result)
placeholder

type opsScheduledReport struct {
	Name       string
	ReportType string
	Schedule   string
	Enabled    bool

	TimeRange time.Duration

	Recipients []string

	ErrorDigestMinCount             int
	AccountHealthErrorRateThreshold float64

	LastRunAt *time.Time
	NextRunAt time.Time
placeholder

type opsScheduledReportContent struct {
	html     string
	overview *OpsDashboardOverview
placeholder

func (s *OpsScheduledReportService) listScheduledReports(ctx context.Context, now time.Time) []*opsScheduledReport {
	if s == nil || s.opsService == nil {
		return nil
placeholder
	if ctx == nil {
		ctx = context.Background()
placeholder

	emailCfg, err := s.opsService.GetEmailNotificationConfig(ctx)
	if err != nil || emailCfg == nil {
		return nil
placeholder
	if !emailCfg.Report.Enabled {
		return nil
placeholder

	recipients := normalizeEmails(emailCfg.Report.Recipients)

	type reportDef struct {
		enabled   bool
		name      string
		kind      string
		timeRange time.Duration
		schedule  string
placeholder

	defs := []reportDef{
		{enabled: emailCfg.Report.DailySummaryEnabled, name: "日报", kind: "daily_summary", timeRange: 24 * time.Hour, schedule: emailCfg.Report.DailySummaryScheduleplaceholder,
		{enabled: emailCfg.Report.WeeklySummaryEnabled, name: "周报", kind: "weekly_summary", timeRange: 7 * 24 * time.Hour, schedule: emailCfg.Report.WeeklySummaryScheduleplaceholder,
		{enabled: emailCfg.Report.ErrorDigestEnabled, name: "错误摘要", kind: "error_digest", timeRange: 24 * time.Hour, schedule: emailCfg.Report.ErrorDigestScheduleplaceholder,
		{enabled: emailCfg.Report.AccountHealthEnabled, name: "账号健康", kind: "account_health", timeRange: 24 * time.Hour, schedule: emailCfg.Report.AccountHealthScheduleplaceholder,
placeholder

	out := make([]*opsScheduledReport, 0, len(defs))
	for _, d := range defs {
		if !d.enabled {
			continue
	placeholder
		spec := strings.TrimSpace(d.schedule)
		if spec == "" {
			continue
	placeholder
		sched, err := opsScheduledReportCronParser.Parse(spec)
		if err != nil {
			log.Printf("[OpsScheduledReport] invalid cron spec=%q for report=%s: %v", spec, d.kind, err)
			continue
	placeholder

		lastRun := s.getLastRunAt(ctx, d.kind)
		base := lastRun
		if base.IsZero() {
			// Allow a schedule matching the current minute to trigger right after startup.
			base = now.Add(-1 * time.Minute)
	placeholder
		next := sched.Next(base)
		if next.IsZero() {
			continue
	placeholder

		var lastRunPtr *time.Time
		if !lastRun.IsZero() {
			lastCopy := lastRun
			lastRunPtr = &lastCopy
	placeholder

		out = append(out, &opsScheduledReport{
			Name:       d.name,
			ReportType: d.kind,
			Schedule:   spec,
			Enabled:    true,

			TimeRange: d.timeRange,

			Recipients: recipients,

			ErrorDigestMinCount:             emailCfg.Report.ErrorDigestMinCount,
			AccountHealthErrorRateThreshold: emailCfg.Report.AccountHealthErrorRateThreshold,

			LastRunAt: lastRunPtr,
			NextRunAt: next,
	placeholder)
placeholder

	return out
placeholder

func (s *OpsScheduledReportService) runReport(ctx context.Context, report *opsScheduledReport, now time.Time) (int, error) {
	if s == nil || s.opsService == nil || s.emailService == nil || report == nil {
		return 0, nil
placeholder
	if ctx == nil {
		ctx = context.Background()
placeholder

	// Mark as "run" up-front so a broken SMTP config doesn't spam retries every minute.
	s.setLastRunAt(ctx, report.ReportType, now)

	content, err := s.generateReportContent(ctx, report, now)
	if err != nil {
		return 0, err
placeholder
	if strings.TrimSpace(content.html) == "" {
		// Skip sending when the report decides not to emit content (e.g., digest below min count).
		return 0, nil
placeholder

	recipients := report.Recipients
	if len(recipients) == 0 && s.userService != nil {
		admin, err := s.userService.GetFirstAdmin(ctx)
		if err == nil && admin != nil && strings.TrimSpace(admin.Email) != "" {
			recipients = []string{strings.TrimSpace(admin.Email)placeholder
	placeholder
placeholder
	if len(recipients) == 0 {
		return 0, nil
placeholder

	attempts := 0
	for _, to := range recipients {
		addr := strings.TrimSpace(to)
		if addr == "" {
			continue
	placeholder
		attempts++
		locale := ""
		if s.emailService.notificationEmailService != nil {
			locale = s.emailService.notificationEmailService.ResolveRecipientLocale(ctx, 0, addr)
			templateVariables := opsScheduledReportLocalizedEmailVariables(report, now, locale)
			rawHTMLVariables := map[string]string{"report_html": content.htmlplaceholder
			if isOpsSummaryReport(report) {
				templateVariables = opsSummaryReportEmailVariables(report, now, content.overview, locale)
		placeholder
			if err := s.emailService.notificationEmailService.Send(ctx, NotificationEmailSendInput{
				Event:            NotificationEmailEventOpsScheduledReport,
				Locale:           locale,
				RecipientEmail:   addr,
				RecipientName:    emailRecipientName(addr),
				SourceType:       "ops_scheduled_report",
				SourceID:         opsScheduledReportDeliverySourceID(report),
				ReminderKey:      now.UTC().Format("2006-01-02T15:04"),
				Variables:        templateVariables,
				RawHTMLVariables: rawHTMLVariables,
		placeholder); err == nil {
				continue
		placeholder else if !shouldFallbackNotificationEmail(err) {
				continue
		placeholder
	placeholder
		subjectName := strings.TrimSpace(report.Name)
		if locale != "" {
			subjectName = opsScheduledReportLocalizedName(report, locale)
	placeholder
		subjectPrefix := "[Ops Report]"
		if strings.HasPrefix(strings.ToLower(locale), "zh") {
			subjectPrefix = "[运维报表]"
	placeholder
		subject := fmt.Sprintf("%s %s", subjectPrefix, subjectName)
		if err := s.emailService.SendEmail(ctx, addr, subject, content.html); err != nil {
			// Ignore per-recipient failures; continue best-effort.
			continue
	placeholder
placeholder
	return attempts, nil
placeholder

func opsScheduledReportDeliverySourceID(report *opsScheduledReport) string {
	if report == nil {
		return "scheduled_report"
placeholder
	parts := []string{
		strings.TrimSpace(report.ReportType),
		strings.TrimSpace(report.Name),
		strings.TrimSpace(report.Schedule),
placeholder
	joined := strings.Trim(strings.Join(parts, ":"), ":")
	if joined == "" {
		return "scheduled_report"
placeholder
	return joined
placeholder

func opsScheduledReportEmailVariables(report *opsScheduledReport, now time.Time) map[string]string {
	end := now.UTC()
	start := end
	name := "Ops report"
	reportType := "scheduled_report"
	if report != nil {
		if strings.TrimSpace(report.Name) != "" {
			name = strings.TrimSpace(report.Name)
	placeholder
		if strings.TrimSpace(report.ReportType) != "" {
			reportType = strings.TrimSpace(report.ReportType)
	placeholder
		if report.TimeRange > 0 {
			start = end.Add(-report.TimeRange)
	placeholder
placeholder
	return map[string]string{
		"report_name":       name,
		"report_type":       reportType,
		"report_start_time": start.Format(time.RFC3339),
		"report_end_time":   end.Format(time.RFC3339),
placeholder
placeholder

func opsScheduledReportLocalizedEmailVariables(report *opsScheduledReport, now time.Time, locale string) map[string]string {
	variables := opsScheduledReportEmailVariables(report, now)
	variables["report_html"] = ""
	variables["report_detail_display"] = "block"
	for _, placeholder := range notificationEmailOpsSummaryPlaceholders {
		variables[placeholder] = "-"
placeholder
	variables["report_summary_display"] = "none"
	if name := opsScheduledReportLocalizedName(report, locale); name != "" {
		variables["report_name"] = name
placeholder
	return variables
placeholder

func opsScheduledReportLocalizedName(report *opsScheduledReport, locale string) string {
	if report == nil {
		return "Ops report"
placeholder
	chinese := strings.HasPrefix(strings.ToLower(strings.TrimSpace(locale)), "zh")
	switch strings.TrimSpace(report.ReportType) {
	case "daily_summary":
		if chinese {
			return "日报"
	placeholder
		return "Daily summary"
	case "weekly_summary":
		if chinese {
			return "周报"
	placeholder
		return "Weekly summary"
	case "error_digest":
		if chinese {
			return "错误摘要"
	placeholder
		return "Error digest"
	case "account_health":
		if chinese {
			return "账号健康"
	placeholder
		return "Account health"
	default:
		return strings.TrimSpace(report.Name)
placeholder
placeholder

func isOpsSummaryReport(report *opsScheduledReport) bool {
	if report == nil {
		return false
placeholder
	switch strings.TrimSpace(report.ReportType) {
	case "daily_summary", "weekly_summary":
		return true
	default:
		return false
placeholder
placeholder

func opsSummaryReportEmailVariables(report *opsScheduledReport, now time.Time, overview *OpsDashboardOverview, locale string) map[string]string {
	variables := opsScheduledReportLocalizedEmailVariables(report, now, locale)
	variables["report_detail_display"] = "none"
	if overview == nil {
		for _, placeholder := range notificationEmailOpsSummaryPlaceholders {
			if placeholder == "report_summary_display" {
				continue
		placeholder
			variables[placeholder] = "-"
	placeholder
		variables["report_summary_display"] = "block"
		return variables
placeholder
	variables["report_summary_display"] = "block"

	variables["report_total_requests"] = formatOpsReportInteger(overview.RequestCountTotal)
	variables["report_success_count"] = formatOpsReportInteger(overview.SuccessCount)
	variables["report_sla_error_count"] = formatOpsReportInteger(overview.ErrorCountSLA)
	variables["report_business_limited_count"] = formatOpsReportInteger(overview.BusinessLimitedCount)
	variables["report_sla"] = fmt.Sprintf("%.2f%%", overview.SLA*100)
	variables["report_error_rate"] = fmt.Sprintf("%.2f%%", overview.ErrorRate*100)
	variables["report_upstream_error_rate"] = fmt.Sprintf("%.2f%%", overview.UpstreamErrorRate*100)
	variables["report_upstream_error_count_excl_429_529"] = formatOpsReportInteger(overview.UpstreamErrorCountExcl429529)
	variables["report_upstream_429_count"] = formatOpsReportInteger(overview.Upstream429Count)
	variables["report_upstream_529_count"] = formatOpsReportInteger(overview.Upstream529Count)
	variables["report_latency_p50"] = formatOpsReportMilliseconds(overview.Duration.P50)
	variables["report_latency_p99"] = formatOpsReportMilliseconds(overview.Duration.P99)
	variables["report_ttft_p50"] = formatOpsReportMilliseconds(overview.TTFT.P50)
	variables["report_ttft_p99"] = formatOpsReportMilliseconds(overview.TTFT.P99)
	variables["report_tokens"] = formatOpsReportInteger(overview.TokenConsumed)
	variables["report_qps_current"] = fmt.Sprintf("%.1f", overview.QPS.Current)
	variables["report_qps_peak"] = fmt.Sprintf("%.1f", overview.QPS.Peak)
	variables["report_qps_avg"] = fmt.Sprintf("%.1f", overview.QPS.Avg)
	variables["report_tps_current"] = fmt.Sprintf("%.1f", overview.TPS.Current)
	variables["report_tps_peak"] = fmt.Sprintf("%.1f", overview.TPS.Peak)
	variables["report_tps_avg"] = fmt.Sprintf("%.1f", overview.TPS.Avg)
	return variables
placeholder

func formatOpsReportInteger(value int64) string {
	raw := strconv.FormatInt(value, 10)
	start := 0
	if strings.HasPrefix(raw, "-") {
		start = 1
placeholder
	if len(raw)-start <= 3 {
		return raw
placeholder

	var builder strings.Builder
	builder.Grow(len(raw) + (len(raw)-start-1)/3)
	_, _ = builder.WriteString(raw[:start])
	digitLen := len(raw) - start
	for offset := 0; offset < digitLen; offset++ {
		if offset > 0 && (digitLen-offset)%3 == 0 {
			_ = builder.WriteByte(',')
	placeholder
		_ = builder.WriteByte(raw[start+offset])
placeholder
	return builder.String()
placeholder

func formatOpsReportMilliseconds(value *int) string {
	if value == nil {
		return "-"
placeholder
	return fmt.Sprintf("%s ms", formatOpsReportInteger(int64(*value)))
placeholder

func (s *OpsScheduledReportService) generateReportContent(ctx context.Context, report *opsScheduledReport, now time.Time) (opsScheduledReportContent, error) {
	if s == nil || s.opsService == nil || report == nil {
		return opsScheduledReportContent{placeholder, fmt.Errorf("service not initialized")
placeholder
	if report.TimeRange <= 0 {
		return opsScheduledReportContent{placeholder, fmt.Errorf("invalid time range")
placeholder

	end := now.UTC()
	start := end.Add(-report.TimeRange)

	switch strings.TrimSpace(report.ReportType) {
	case "daily_summary", "weekly_summary":
		overview, err := s.opsService.GetDashboardOverview(ctx, &OpsDashboardFilter{
			StartTime: start,
			EndTime:   end,
			Platform:  "",
			GroupID:   nil,
			QueryMode: OpsQueryModeAuto,
	placeholder)
		if err != nil {
			// If pre-aggregation isn't ready but the report is requested, fall back to raw.
			if strings.TrimSpace(report.ReportType) == "daily_summary" || strings.TrimSpace(report.ReportType) == "weekly_summary" {
				overview, err = s.opsService.GetDashboardOverview(ctx, &OpsDashboardFilter{
					StartTime: start,
					EndTime:   end,
					Platform:  "",
					GroupID:   nil,
					QueryMode: OpsQueryModeRaw,
			placeholder)
		placeholder
			if err != nil {
				return opsScheduledReportContent{placeholder, err
		placeholder
	placeholder
		return opsScheduledReportContent{
			html:     buildOpsSummaryEmailHTML(report.Name, start, end, overview),
			overview: overview,
	placeholder, nil
	case "error_digest":
		// Lightweight digest: list recent errors (status>=400) and breakdown by type.
		startTime := start
		endTime := end
		filter := &OpsErrorLogFilter{
			StartTime: &startTime,
			EndTime:   &endTime,
			Page:      1,
			PageSize:  100,
	placeholder
		out, err := s.opsService.GetErrorLogs(ctx, filter)
		if err != nil {
			return opsScheduledReportContent{placeholder, err
	placeholder
		if report.ErrorDigestMinCount > 0 && out != nil && out.Total < report.ErrorDigestMinCount {
			return opsScheduledReportContent{placeholder, nil
	placeholder
		return opsScheduledReportContent{html: buildOpsErrorDigestEmailHTML(report.Name, start, end, out)placeholder, nil
	case "account_health":
		// Best-effort: use account availability (not error rate yet).
		avail, err := s.opsService.GetAccountAvailability(ctx, "", nil)
		if err != nil {
			return opsScheduledReportContent{placeholder, err
	placeholder
		_ = report.AccountHealthErrorRateThreshold // reserved for future per-account error rate report
		return opsScheduledReportContent{html: buildOpsAccountHealthEmailHTML(report.Name, start, end, avail)placeholder, nil
	default:
		return opsScheduledReportContent{placeholder, fmt.Errorf("unknown report type: %s", report.ReportType)
placeholder
placeholder

func buildOpsSummaryEmailHTML(title string, start, end time.Time, overview *OpsDashboardOverview) string {
	if overview == nil {
		return fmt.Sprintf("<h2>%s</h2><p>No data.</p>", htmlEscape(title))
placeholder

	latP50 := "-"
	latP99 := "-"
	if overview.Duration.P50 != nil {
		latP50 = fmt.Sprintf("%dms", *overview.Duration.P50)
placeholder
	if overview.Duration.P99 != nil {
		latP99 = fmt.Sprintf("%dms", *overview.Duration.P99)
placeholder

	ttftP50 := "-"
	ttftP99 := "-"
	if overview.TTFT.P50 != nil {
		ttftP50 = fmt.Sprintf("%dms", *overview.TTFT.P50)
placeholder
	if overview.TTFT.P99 != nil {
		ttftP99 = fmt.Sprintf("%dms", *overview.TTFT.P99)
placeholder

	return fmt.Sprintf(`
<h2>%s</h2>
<p><b>Period</b>: %s ~ %s (UTC)</p>
<ul>
  <li><b>Total Requests</b>: %d</li>
  <li><b>Success</b>: %d</li>
  <li><b>Errors (SLA)</b>: %d</li>
  <li><b>Business Limited</b>: %d</li>
  <li><b>SLA</b>: %.2f%%</li>
  <li><b>Error Rate</b>: %.2f%%</li>
  <li><b>Upstream Error Rate (excl 429/529)</b>: %.2f%%</li>
  <li><b>Upstream Errors</b>: excl429/529=%d, 429=%d, 529=%d</li>
  <li><b>Latency</b>: p50=%s, p99=%s</li>
  <li><b>TTFT</b>: p50=%s, p99=%s</li>
  <li><b>Tokens</b>: %d</li>
  <li><b>QPS</b>: current=%.1f, peak=%.1f, avg=%.1f</li>
  <li><b>TPS</b>: current=%.1f, peak=%.1f, avg=%.1f</li>
</ul>
`,
		htmlEscape(strings.TrimSpace(title)),
		htmlEscape(start.UTC().Format(time.RFC3339)),
		htmlEscape(end.UTC().Format(time.RFC3339)),
		overview.RequestCountTotal,
		overview.SuccessCount,
		overview.ErrorCountSLA,
		overview.BusinessLimitedCount,
		overview.SLA*100,
		overview.ErrorRate*100,
		overview.UpstreamErrorRate*100,
		overview.UpstreamErrorCountExcl429529,
		overview.Upstream429Count,
		overview.Upstream529Count,
		htmlEscape(latP50),
		htmlEscape(latP99),
		htmlEscape(ttftP50),
		htmlEscape(ttftP99),
		overview.TokenConsumed,
		overview.QPS.Current,
		overview.QPS.Peak,
		overview.QPS.Avg,
		overview.TPS.Current,
		overview.TPS.Peak,
		overview.TPS.Avg,
	)
placeholder

func buildOpsErrorDigestEmailHTML(title string, start, end time.Time, list *OpsErrorLogList) string {
	total := 0
	recent := []*OpsErrorLog{placeholder
	if list != nil {
		total = list.Total
		recent = list.Errors
placeholder
	if len(recent) > 10 {
		recent = recent[:10]
placeholder

	rows := ""
	for _, item := range recent {
		if item == nil {
			continue
	placeholder
		rows += fmt.Sprintf(
			"<tr><td>%s</td><td>%s</td><td>%d</td><td>%s</td></tr>",
			htmlEscape(item.CreatedAt.UTC().Format(time.RFC3339)),
			htmlEscape(item.Platform),
			item.StatusCode,
			htmlEscape(truncateString(item.Message, 180)),
		)
placeholder
	if rows == "" {
		rows = "<tr><td colspan=\"4\">No recent errors.</td></tr>"
placeholder

	return fmt.Sprintf(`
<h2>%s</h2>
<p><b>Period</b>: %s ~ %s (UTC)</p>
<p><b>Total Errors</b>: %d</p>
<h3>Recent</h3>
<table border="1" cellpadding="6" cellspacing="0" style="border-collapse:collapse;">
  <thead><tr><th>Time</th><th>Platform</th><th>Status</th><th>Message</th></tr></thead>
  <tbody>%s</tbody>
</table>
`,
		htmlEscape(strings.TrimSpace(title)),
		htmlEscape(start.UTC().Format(time.RFC3339)),
		htmlEscape(end.UTC().Format(time.RFC3339)),
		total,
		rows,
	)
placeholder

func buildOpsAccountHealthEmailHTML(title string, start, end time.Time, avail *OpsAccountAvailability) string {
	total := 0
	available := 0
	rateLimited := 0
	hasError := 0

	if avail != nil && avail.Accounts != nil {
		for _, a := range avail.Accounts {
			if a == nil {
				continue
		placeholder
			total++
			if a.IsAvailable {
				available++
		placeholder
			if a.IsRateLimited {
				rateLimited++
		placeholder
			if a.HasError {
				hasError++
		placeholder
	placeholder
placeholder

	return fmt.Sprintf(`
<h2>%s</h2>
<p><b>Period</b>: %s ~ %s (UTC)</p>
<ul>
  <li><b>Total Accounts</b>: %d</li>
  <li><b>Available</b>: %d</li>
  <li><b>Rate Limited</b>: %d</li>
  <li><b>Error</b>: %d</li>
</ul>
<p>Note: This report currently reflects account availability status only.</p>
`,
		htmlEscape(strings.TrimSpace(title)),
		htmlEscape(start.UTC().Format(time.RFC3339)),
		htmlEscape(end.UTC().Format(time.RFC3339)),
		total,
		available,
		rateLimited,
		hasError,
	)
placeholder

func (s *OpsScheduledReportService) tryAcquireLeaderLock(ctx context.Context) (func(), bool) {
	if s == nil || !s.distributedLockOn {
		return nil, true
placeholder
	if s.redisClient == nil {
		s.warnNoRedisOnce.Do(func() {
			log.Printf("[OpsScheduledReport] redis not configured; running without distributed lock")
	placeholder)
		return nil, true
placeholder
	if ctx == nil {
		ctx = context.Background()
placeholder

	key := opsScheduledReportLeaderLockKeyDefault
	ttl := opsScheduledReportLeaderLockTTLDefault
	if strings.TrimSpace(key) == "" {
		key = "ops:scheduled_reports:leader"
placeholder
	if ttl <= 0 {
		ttl = 5 * time.Minute
placeholder

	ok, err := s.redisClient.SetNX(ctx, key, s.instanceID, ttl).Result()
	if err != nil {
		// Prefer fail-closed to avoid duplicate report sends when Redis is flaky.
		log.Printf("[OpsScheduledReport] leader lock SetNX failed; skipping this cycle: %v", err)
		return nil, false
placeholder
	if !ok {
		return nil, false
placeholder
	return func() {
		_, _ = opsScheduledReportReleaseScript.Run(ctx, s.redisClient, []string{keyplaceholder, s.instanceID).Result()
placeholder, true
placeholder

func (s *OpsScheduledReportService) getLastRunAt(ctx context.Context, reportType string) time.Time {
	if s == nil || s.redisClient == nil {
		return time.Time{placeholder
placeholder
	kind := strings.TrimSpace(reportType)
	if kind == "" {
		return time.Time{placeholder
placeholder
	key := opsScheduledReportLastRunKeyPrefix + kind

	raw, err := s.redisClient.Get(ctx, key).Result()
	if err != nil || strings.TrimSpace(raw) == "" {
		return time.Time{placeholder
placeholder
	sec, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || sec <= 0 {
		return time.Time{placeholder
placeholder
	last := time.Unix(sec, 0)
	// Cron schedules are interpreted in the configured timezone (s.loc). Ensure the base time
	// passed into cron.Next() uses the same location; otherwise the job will drift by timezone
	// offset (e.g. Asia/Shanghai default would run 8h later after the first execution).
	if s.loc != nil {
		return last.In(s.loc)
placeholder
	return last.UTC()
placeholder

func (s *OpsScheduledReportService) setLastRunAt(ctx context.Context, reportType string, t time.Time) {
	if s == nil || s.redisClient == nil {
		return
placeholder
	kind := strings.TrimSpace(reportType)
	if kind == "" {
		return
placeholder
	if t.IsZero() {
		t = time.Now().UTC()
placeholder
	key := opsScheduledReportLastRunKeyPrefix + kind
	_ = s.redisClient.Set(ctx, key, strconv.FormatInt(t.UTC().Unix(), 10), 14*24*time.Hour).Err()
placeholder

func (s *OpsScheduledReportService) recordHeartbeatSuccess(runAt time.Time, duration time.Duration, result string) {
	if s == nil || s.opsService == nil || s.opsService.opsRepo == nil {
		return
placeholder
	now := time.Now().UTC()
	durMs := duration.Milliseconds()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := strings.TrimSpace(result)
	if msg == "" {
		msg = "ok"
placeholder
	msg = truncateString(msg, 2048)
	_ = s.opsService.opsRepo.UpsertJobHeartbeat(ctx, &OpsUpsertJobHeartbeatInput{
		JobName:        opsScheduledReportJobName,
		LastRunAt:      &runAt,
		LastSuccessAt:  &now,
		LastDurationMs: &durMs,
		LastResult:     &msg,
placeholder)
placeholder

func (s *OpsScheduledReportService) recordHeartbeatError(runAt time.Time, duration time.Duration, err error) {
	if s == nil || s.opsService == nil || s.opsService.opsRepo == nil || err == nil {
		return
placeholder
	now := time.Now().UTC()
	durMs := duration.Milliseconds()
	msg := truncateString(err.Error(), 2048)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = s.opsService.opsRepo.UpsertJobHeartbeat(ctx, &OpsUpsertJobHeartbeatInput{
		JobName:        opsScheduledReportJobName,
		LastRunAt:      &runAt,
		LastErrorAt:    &now,
		LastError:      &msg,
		LastDurationMs: &durMs,
placeholder)
placeholder

func normalizeEmails(in []string) []string {
	if len(in) == 0 {
		return nil
placeholder
	seen := make(map[string]struct{placeholder, len(in))
	out := make([]string, 0, len(in))
	for _, raw := range in {
		addr := strings.ToLower(strings.TrimSpace(raw))
		if addr == "" {
			continue
	placeholder
		if _, ok := seen[addr]; ok {
			continue
	placeholder
		seen[addr] = struct{placeholder{placeholder
		out = append(out, addr)
placeholder
	return out
placeholder
