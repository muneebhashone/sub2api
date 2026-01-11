//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/repository"
	"github.com/Wei-Shaw/sub2api/internal/server"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Server  *http.Server
	Cleanup func()
placeholder

func initializeApplication(buildInfo handler.BuildInfo) (*Application, error) {
	wire.Build(
		// Infrastructure layer ProviderSets
		config.ProviderSet,

		// Business layer ProviderSets
		repository.ProviderSet,
		service.ProviderSet,
		middleware.ProviderSet,
		handler.ProviderSet,

		// Server layer ProviderSet
		server.ProviderSet,

		// BuildInfo provider
		provideServiceBuildInfo,

		// Cleanup function provider
		provideCleanup,

		// Application struct
		wire.Struct(new(Application), "Server", "Cleanup"),
	)
	return nil, nil
placeholder

func provideServiceBuildInfo(buildInfo handler.BuildInfo) service.BuildInfo {
	return service.BuildInfo{
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
	tokenRefresh *service.TokenRefreshService,
	accountExpiry *service.AccountExpiryService,
	pricing *service.PricingService,
	emailQueue *service.EmailQueueService,
	billingCache *service.BillingCacheService,
	oauth *service.OAuthService,
	openaiOAuth *service.OpenAIOAuthService,
	geminiOAuth *service.GeminiOAuthService,
	antigravityOAuth *service.AntigravityOAuthService,
) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Cleanup steps in reverse dependency order
		cleanupSteps := []struct {
			name string
			fn   func() error
	placeholder{
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
			{"TokenRefreshService", func() error {
				tokenRefresh.Stop()
				return nil
	placeholder
			{"AccountExpiryService", func() error {
				accountExpiry.Stop()
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
			{"Redis", func() error {
				return rdb.Close()
	placeholder
			{"Ent", func() error {
				return entClient.Close()
	placeholder
	placeholder

		for _, step := range cleanupSteps {
			if err := step.fn(); err != nil {
				log.Printf("[Cleanup] %s failed: %v", step.name, err)
				// Continue with remaining cleanup steps even if one fails
		placeholder else {
				log.Printf("[Cleanup] %s succeeded", step.name)
		placeholder
	placeholder

		// Check if context timed out
		select {
		case <-ctx.Done():
			log.Printf("[Cleanup] Warning: cleanup timed out after 10 seconds")
		default:
			log.Printf("[Cleanup] All cleanup steps completed")
	placeholder
placeholder
placeholder
