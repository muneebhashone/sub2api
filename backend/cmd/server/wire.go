//go:build wireinject
// +build wireinject

package main

import (
	"sub2api/internal/config"
	"sub2api/internal/handler"
	"sub2api/internal/infrastructure"
	"sub2api/internal/repository"
	"sub2api/internal/server"
	"sub2api/internal/service"

	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	Server  *http.Server
	Cleanup func()
placeholder

func initializeApplication(buildInfo handler.BuildInfo) (*Application, error) {
	wire.Build(
		// 基础设施层 ProviderSets
		config.ProviderSet,
		infrastructure.ProviderSet,

		// 业务层 ProviderSets
		repository.ProviderSet,
		service.ProviderSet,
		handler.ProviderSet,

		// 服务器层 ProviderSet
		server.ProviderSet,

		// BuildInfo provider
		provideServiceBuildInfo,

		// 清理函数提供者
		provideCleanup,

		// 应用程序结构体
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
	db *gorm.DB,
	rdb *redis.Client,
	services *service.Services,
) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Cleanup steps in reverse dependency order
		cleanupSteps := []struct {
			name string
			fn   func() error
	placeholder{
			{"TokenRefreshService", func() error {
				services.TokenRefresh.Stop()
				return nil
	placeholder
			{"PricingService", func() error {
				services.Pricing.Stop()
				return nil
	placeholder
			{"EmailQueueService", func() error {
				services.EmailQueue.Stop()
				return nil
	placeholder
			{"Redis", func() error {
				return rdb.Close()
	placeholder
			{"Database", func() error {
				sqlDB, err := db.DB()
				if err != nil {
					return err
			placeholder
				return sqlDB.Close()
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
