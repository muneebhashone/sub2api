//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/Wei-Shaw/sub2api/internal/repository"
	"github.com/Wei-Shaw/sub2api/internal/securityaudit"
	"github.com/Wei-Shaw/sub2api/internal/server"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Server        *http.Server
	PromptAudit   *securityaudit.PromptService
	PluginManager *service.PluginManager
	Cleanup       func()
placeholder

func initializeApplication(buildInfo handler.BuildInfo) (*Application, error) {
	wire.Build(
		// Infrastructure layer ProviderSets
		config.ProviderSet,

		// Business layer ProviderSets
		repository.ProviderSet,
		service.ProviderSet,
		securityaudit.ProviderSet,
		payment.ProviderSet,
		middleware.ProviderSet,
		handler.ProviderSet,

		// Server layer ProviderSet
		server.ProviderSet,

		// Privacy client factory for OpenAI training opt-out
		providePrivacyClientFactory,

		// BuildInfo provider
		provideServiceBuildInfo,
		providePluginHostInfo,

		// Cleanup function provider
		provideCleanup,

		// Application struct
		wire.Struct(new(Application), "Server", "PromptAudit", "PluginManager", "Cleanup"),
	)
	return nil, nil
placeholder

func providePrivacyClientFactory() service.PrivacyClientFactory {
	return repository.CreatePrivacyReqClient
placeholder

func provideServiceBuildInfo(buildInfo handler.BuildInfo) service.BuildInfo {
	return service.BuildInfo{
		Version:   buildInfo.Version,
		BuildType: buildInfo.BuildType,
placeholder
placeholder

func providePluginHostInfo(buildInfo handler.BuildInfo) service.PluginHostInfo {
	return service.PluginHostInfo{
		Version:   buildInfo.Version,
		BuildType: buildInfo.BuildType,
placeholder
placeholder

func provideCleanup(
	entClient *ent.Client,
	rdb *redis.Client,
	opsMetricsCollector *service.OpsMetricsCollector,
	opsAggregation *service.OpsAggregationService,
	opsAlertEvaluator *service.OpsAlertEvaluatorService,
	opsCleanup *service.OpsCleanupService,
	opsScheduledReport *service.OpsScheduledReportService,
	opsSystemLogSink *service.OpsSystemLogSink,
	opsService *service.OpsService,
	opsIngressReject *service.OpsIngressRejectAggregator,
	apiKeyService *service.APIKeyService,
	authCacheInvalidationWorker *service.AuthCacheInvalidationWorker,
	schedulerSnapshot *service.SchedulerSnapshotService,
	tokenRefresh *service.TokenRefreshService,
	accountExpiry *service.AccountExpiryService,
	cnProviderBalanceCheck *service.CNProviderBalanceCheckService,
	codexVersionSync *service.OpenAICodexVersionSyncService,
	proxyExpiry *service.ProxyExpiryService,
	subscriptionExpiry *service.SubscriptionExpiryService,
	usageCleanup *service.UsageCleanupService,
	idempotencyCleanup *service.IdempotencyCleanupService,
	batchImageCleanup *service.BatchImageCleanupService,
	batchImageWorker *service.BatchImageWorkerRuntime,
	pricing *service.PricingService,
	emailQueue *service.EmailQueueService,
	billingCache *service.BillingCacheService,
	usageRecordWorkerPool *service.UsageRecordWorkerPool,
	subscriptionService *service.SubscriptionService,
	oauth *service.OAuthService,
	openaiOAuth *service.OpenAIOAuthService,
	geminiOAuth *service.GeminiOAuthService,
	antigravityOAuth *service.AntigravityOAuthService,
	grokOAuth *service.GrokOAuthService,
	openAIGateway *service.OpenAIGatewayService,
	scheduledTestRunner *service.ScheduledTestRunnerService,
	backupSvc *service.BackupService,
	paymentOrderExpiry *service.PaymentOrderExpiryService,
	channelMonitorRunner *service.ChannelMonitorRunner,
	channelMonitorV2Aggregator *service.ChannelMonitorV2Aggregator,
	quotaFlusher *service.UserPlatformQuotaUsageFlusher,
	upstreamBillingProbe *service.UpstreamBillingProbeService,
	ollamaCloudUsage *service.OllamaCloudUsageService,
	auditLog *service.AuditLogService,
	promptAudit *securityaudit.PromptService,
	pluginManager *service.PluginManager,
) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		type cleanupStep struct {
			name string
			fn   func() error
	placeholder

		// 应用层清理步骤可并行执行，基础设施资源（Redis/Ent）最后按顺序关闭。
		parallelSteps := []cleanupStep{
			{"PluginManager", func() error {
				if pluginManager != nil {
					pluginManager.Stop()
			placeholder
				return nil
	placeholder
			{"OpsIngressRejectAggregator", func() error {
				if opsIngressReject != nil {
					opsIngressReject.Stop()
			placeholder
				return nil
	placeholder
			{"AuthCacheInvalidationWorker", func() error {
				if authCacheInvalidationWorker != nil {
					authCacheInvalidationWorker.Stop()
			placeholder
				return nil
	placeholder
			{"AuthCacheInvalidationSubscriber", func() error {
				if apiKeyService != nil {
					apiKeyService.StopAuthCacheInvalidationSubscriber()
			placeholder
				return nil
	placeholder
			{"OpsRuntimeSettingsRefresh", func() error {
				if opsService != nil {
					opsService.StopRuntimeSettingsRefresh()
			placeholder
				return nil
	placeholder
			{"PromptAuditService", func() error {
				if promptAudit != nil {
					return promptAudit.Shutdown(ctx)
			placeholder
				return nil
	placeholder
			{"OpsScheduledReportService", func() error {
				if opsScheduledReport != nil {
					opsScheduledReport.Stop()
			placeholder
				return nil
	placeholder
			{"OpsCleanupService", func() error {
				if opsCleanup != nil {
					opsCleanup.Stop()
			placeholder
				return nil
	placeholder
			{"OpsSystemLogSink", func() error {
				if opsSystemLogSink != nil {
					opsSystemLogSink.Stop()
			placeholder
				return nil
	placeholder
			{"AuditLogService", func() error {
				if auditLog != nil {
					auditLog.Stop()
			placeholder
				return nil
	placeholder
			{"OpsAlertEvaluatorService", func() error {
				if opsAlertEvaluator != nil {
					opsAlertEvaluator.Stop()
			placeholder
				return nil
	placeholder
			{"OpsAggregationService", func() error {
				if opsAggregation != nil {
					opsAggregation.Stop()
			placeholder
				return nil
	placeholder
			{"OpsMetricsCollector", func() error {
				if opsMetricsCollector != nil {
					opsMetricsCollector.Stop()
			placeholder
				return nil
	placeholder
			{"SchedulerSnapshotService", func() error {
				if schedulerSnapshot != nil {
					schedulerSnapshot.Stop()
			placeholder
				return nil
	placeholder
			{"UsageCleanupService", func() error {
				if usageCleanup != nil {
					usageCleanup.Stop()
			placeholder
				return nil
	placeholder
			{"IdempotencyCleanupService", func() error {
				if idempotencyCleanup != nil {
					idempotencyCleanup.Stop()
			placeholder
				return nil
	placeholder
			{"BatchImageCleanupService", func() error {
				if batchImageCleanup != nil {
					batchImageCleanup.Stop()
			placeholder
				return nil
	placeholder
			{"BatchImageWorkerRuntime", func() error {
				if batchImageWorker != nil {
					batchImageWorker.Stop()
			placeholder
				return nil
	placeholder
			{"TokenRefreshService", func() error {
				tokenRefresh.Stop()
				return nil
	placeholder
			{"AccountExpiryService", func() error {
				accountExpiry.Stop()
				return nil
	placeholder
			{"CNProviderBalanceCheckService", func() error {
				if cnProviderBalanceCheck != nil {
					cnProviderBalanceCheck.Stop()
			placeholder
				return nil
	placeholder
			{"OpenAICodexVersionSyncService", func() error {
				codexVersionSync.Stop()
				return nil
	placeholder
			{"ProxyExpiryService", func() error {
				proxyExpiry.Stop()
				return nil
	placeholder
			{"SubscriptionExpiryService", func() error {
				subscriptionExpiry.Stop()
				return nil
	placeholder
			{"SubscriptionService", func() error {
				if subscriptionService != nil {
					subscriptionService.Stop()
			placeholder
				return nil
	placeholder
			{"PricingService", func() error {
				pricing.Stop()
				return nil
	placeholder
			{"EmailQueueService", func() error {
				emailQueue.Stop()
				return nil
	placeholder
			{"BillingCacheService", func() error {
				billingCache.Stop()
				return nil
	placeholder
			{"UsageRecordWorkerPool", func() error {
				if usageRecordWorkerPool != nil {
					usageRecordWorkerPool.Stop()
			placeholder
				return nil
	placeholder
			{"OAuthService", func() error {
				oauth.Stop()
				return nil
	placeholder
			{"OpenAIOAuthService", func() error {
				openaiOAuth.Stop()
				return nil
	placeholder
			{"GeminiOAuthService", func() error {
				geminiOAuth.Stop()
				return nil
	placeholder
			{"AntigravityOAuthService", func() error {
				antigravityOAuth.Stop()
				return nil
	placeholder
			{"GrokOAuthService", func() error {
				if grokOAuth != nil {
					grokOAuth.Stop()
			placeholder
				return nil
	placeholder
			{"OpenAIWSPool", func() error {
				if openAIGateway != nil {
					openAIGateway.CloseOpenAIWSPool()
			placeholder
				return nil
	placeholder
			{"ScheduledTestRunnerService", func() error {
				if scheduledTestRunner != nil {
					scheduledTestRunner.Stop()
			placeholder
				return nil
	placeholder
			{"BackupService", func() error {
				if backupSvc != nil {
					backupSvc.Stop()
			placeholder
				return nil
	placeholder
			{"PaymentOrderExpiryService", func() error {
				if paymentOrderExpiry != nil {
					paymentOrderExpiry.Stop()
			placeholder
				return nil
	placeholder
			{"ChannelMonitorV2Aggregator", func() error {
				if channelMonitorV2Aggregator != nil {
					channelMonitorV2Aggregator.Stop()
			placeholder
				return nil
	placeholder
			{"ChannelMonitorRunner", func() error {
				if channelMonitorRunner != nil {
					channelMonitorRunner.Stop()
			placeholder
				return nil
	placeholder
			{"UserPlatformQuotaUsageFlusher", func() error {
				if quotaFlusher != nil {
					quotaFlusher.Stop()
			placeholder
				return nil
	placeholder
			{"UpstreamBillingProbeService", func() error {
				if upstreamBillingProbe != nil {
					upstreamBillingProbe.Stop()
			placeholder
				return nil
	placeholder
			{"OllamaCloudUsageService", func() error {
				if ollamaCloudUsage != nil {
					ollamaCloudUsage.Stop()
			placeholder
				return nil
	placeholder
	placeholder

		infraSteps := []cleanupStep{
			{"Redis", func() error {
				if rdb == nil {
					return nil
			placeholder
				return rdb.Close()
	placeholder
			{"Ent", func() error {
				if entClient == nil {
					return nil
			placeholder
				return entClient.Close()
	placeholder
	placeholder

		runParallel := func(steps []cleanupStep) {
			var wg sync.WaitGroup
			for i := range steps {
				step := steps[i]
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := step.fn(); err != nil {
						log.Printf("[Cleanup] %s failed: %v", step.name, err)
						return
				placeholder
					log.Printf("[Cleanup] %s succeeded", step.name)
			placeholder()
		placeholder
			wg.Wait()
	placeholder

		runSequential := func(steps []cleanupStep) {
			for i := range steps {
				step := steps[i]
				if err := step.fn(); err != nil {
					log.Printf("[Cleanup] %s failed: %v", step.name, err)
					continue
			placeholder
				log.Printf("[Cleanup] %s succeeded", step.name)
		placeholder
	placeholder

		runParallel(parallelSteps)
		runSequential(infraSteps)

		// Check if context timed out
		select {
		case <-ctx.Done():
			log.Printf("[Cleanup] Warning: cleanup timed out after 10 seconds")
		default:
			log.Printf("[Cleanup] All cleanup steps completed")
	placeholder
placeholder
placeholder
