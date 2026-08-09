package admin

import (
	"context"
	"log/slog"
	"sync"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

const (
	grokImportProbeConcurrency = 3
	grokImportProbeTimeout     = 25 * time.Second
	grokImportProbeQueueLimit  = 64
)

type grokImportProber interface {
	QueryQuota(ctx context.Context, accountID int64) (*service.GrokQuotaProbeResult, error)
placeholder

type grokImportProbeTask struct {
	prober    grokImportProber
	accountID int64
placeholder

type grokImportProbeScheduler struct {
	mu          sync.Mutex
	queue       []grokImportProbeTask
	pending     map[int64]struct{placeholder
	inFlight    map[int64]struct{placeholder
	concurrency int
	workers     int
	maxWorkers  int
	timeout     time.Duration
placeholder

var defaultGrokImportProbeScheduler = newGrokImportProbeScheduler(
	grokImportProbeConcurrency,
	grokImportProbeTimeout,
)

func newGrokImportProbeScheduler(concurrency int, timeout time.Duration) *grokImportProbeScheduler {
	if concurrency <= 0 {
		concurrency = 1
placeholder
	if timeout <= 0 {
		timeout = grokImportProbeTimeout
placeholder
	return &grokImportProbeScheduler{
		concurrency: concurrency,
		timeout:     timeout,
		pending:     make(map[int64]struct{placeholder),
		inFlight:    make(map[int64]struct{placeholder),
placeholder
placeholder

func (s *grokImportProbeScheduler) schedule(prober grokImportProber, account *service.Account) {
	if s == nil || prober == nil || account == nil || account.ID <= 0 {
		return
placeholder
	if account.Platform != service.PlatformGrok || account.Type != service.AccountTypeOAuth {
		return
placeholder

	s.mu.Lock()
	if _, exists := s.pending[account.ID]; exists {
		s.mu.Unlock()
		return
placeholder
	if _, exists := s.inFlight[account.ID]; exists {
		s.mu.Unlock()
		return
placeholder
	if len(s.queue) >= grokImportProbeQueueLimit {
		s.mu.Unlock()
		slog.Debug("grok_import_active_probe_dropped", "account_id", account.ID, "reason", "queue_full")
		return
placeholder
	s.queue = append(s.queue, grokImportProbeTask{prober: prober, accountID: account.IDplaceholder)
	s.pending[account.ID] = struct{placeholder{placeholder
	if s.workers < s.concurrency {
		s.workers++
		if s.workers > s.maxWorkers {
			s.maxWorkers = s.workers
	placeholder
		go s.worker()
placeholder
	s.mu.Unlock()
placeholder

func (s *grokImportProbeScheduler) worker() {
	for {
		task, ok := s.nextTask()
		if !ok {
			return
	placeholder
		s.run(task.prober, task.accountID)
		s.finish(task.accountID)
placeholder
placeholder

func (s *grokImportProbeScheduler) nextTask() (grokImportProbeTask, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.queue) == 0 {
		s.workers--
		return grokImportProbeTask{placeholder, false
placeholder
	task := s.queue[0]
	s.queue[0] = grokImportProbeTask{placeholder
	s.queue = s.queue[1:]
	if len(s.queue) == 0 {
		s.queue = nil
placeholder
	delete(s.pending, task.accountID)
	s.inFlight[task.accountID] = struct{placeholder{placeholder
	return task, true
placeholder

func (s *grokImportProbeScheduler) finish(accountID int64) {
	s.mu.Lock()
	delete(s.inFlight, accountID)
	s.mu.Unlock()
placeholder

func (s *grokImportProbeScheduler) run(prober grokImportProber, accountID int64) {
	defer func() {
		if recovered := recover(); recovered != nil {
			slog.Error(
				"grok_import_active_probe_panic",
				"account_id", accountID,
				"recovery_type", panicType(recovered),
			)
	placeholder
placeholder()

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	result, err := prober.QueryQuota(ctx, accountID)
	if err != nil {
		slog.Warn(
			"grok_import_active_probe_failed",
			"account_id", accountID,
			"status", infraerrors.Code(err),
			"reason", infraerrors.Reason(err),
		)
		return
placeholder
	if result == nil {
		slog.Warn(
			"grok_import_active_probe_failed",
			"account_id", accountID,
			"reason", "empty_result",
		)
		return
placeholder

	slog.Info(
		"grok_import_active_probe_completed",
		"account_id", accountID,
		"model", result.Model,
		"status", result.StatusCode,
		"headers_observed", result.HeadersObserved,
	)
placeholder

func panicType(value any) string {
	switch value.(type) {
	case string:
		return "string"
	case error:
		return "error"
	default:
		return "unknown"
placeholder
placeholder

func (h *AccountHandler) scheduleGrokImportProbe(account *service.Account) {
	if h == nil {
		return
placeholder
	defaultGrokImportProbeScheduler.schedule(h.grokImportProber, account)
placeholder

func (h *GrokOAuthHandler) scheduleGrokImportProbe(account *service.Account) {
	if h == nil {
		return
placeholder
	defaultGrokImportProbeScheduler.schedule(h.importProber, account)
placeholder

// ProvideAccountHandler injects the Grok active prober for production while
// keeping NewAccountHandler convenient for focused unit tests.
func ProvideAccountHandler(
	adminService service.AdminService,
	oauthService *service.OAuthService,
	openaiOAuthService *service.OpenAIOAuthService,
	geminiOAuthService *service.GeminiOAuthService,
	antigravityOAuthService *service.AntigravityOAuthService,
	grokOAuthService service.GrokOAuthTokenService,
	rateLimitService *service.RateLimitService,
	accountUsageService *service.AccountUsageService,
	accountTestService *service.AccountTestService,
	concurrencyService *service.ConcurrencyService,
	crsSyncService *service.CRSSyncService,
	sessionLimitCache service.SessionLimitCache,
	rpmCache service.RPMCache,
	tokenCacheInvalidator service.TokenCacheInvalidator,
	grokQuotaService *service.GrokQuotaService,
) *AccountHandler {
	handler := NewAccountHandler(
		adminService,
		oauthService,
		openaiOAuthService,
		geminiOAuthService,
		antigravityOAuthService,
		grokOAuthService,
		rateLimitService,
		accountUsageService,
		accountTestService,
		concurrencyService,
		crsSyncService,
		sessionLimitCache,
		rpmCache,
		tokenCacheInvalidator,
	)
	handler.grokImportProber = grokQuotaService
	return handler
placeholder
