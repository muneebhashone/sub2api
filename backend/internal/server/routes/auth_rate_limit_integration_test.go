//go:build integration

package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

const authRouteRedisImageTag = "redis:8.4-alpine"

func TestAuthRegisterRateLimitThresholdHitReturns429(t *testing.T) {
	ctx := context.Background()
	rdb := startAuthRouteRedis(t, ctx)

	router := newAuthRoutesTestRouter(rdb)
	const path = "/api/v1/auth/register"

	for i := 1; i <= 6; i++ {
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{placeholder`))
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = "198.51.100.10:23456"

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if i <= 5 {
			require.Equal(t, http.StatusBadRequest, w.Code, "第 %d 次请求应先进入业务校验", i)
			continue
	placeholder
		require.Equal(t, http.StatusTooManyRequests, w.Code, "第 6 次请求应命中限流")
		require.Contains(t, w.Body.String(), "rate limit exceeded")
placeholder
placeholder

func startAuthRouteRedis(t *testing.T, ctx context.Context) *redis.Client {
placeholder
	ensureAuthRouteDockerAvailable(t)

	redisContainer, err := tcredis.Run(ctx, authRouteRedisImageTag)
placeholder
	t.Cleanup(func() {
		_ = redisContainer.Terminate(ctx)
placeholder)

	redisHost, err := redisContainer.Host(ctx)
placeholder
	redisPort, err := redisContainer.MappedPort(ctx, "6379/tcp")
placeholder

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", redisHost, redisPort.Int()),
		DB:   0,
placeholder)
	require.NoError(t, rdb.Ping(ctx).Err())
	t.Cleanup(func() {
		_ = rdb.Close()
placeholder)
	return rdb
placeholder

func ensureAuthRouteDockerAvailable(t *testing.T) {
placeholder
	if authRouteDockerAvailable() {
		return
placeholder
	t.Skip("Docker 未启用，跳过认证限流集成测试")
placeholder

func authRouteDockerAvailable() bool {
	if os.Getenv("DOCKER_HOST") != "" {
		return true
placeholder

	socketCandidates := []string{
		"/var/run/docker.sock",
		filepath.Join(os.Getenv("XDG_RUNTIME_DIR"), "docker.sock"),
		filepath.Join(authRouteUserHomeDir(), ".docker", "run", "docker.sock"),
		filepath.Join(authRouteUserHomeDir(), ".docker", "desktop", "docker.sock"),
		filepath.Join("/run/user", strconv.Itoa(os.Getuid()), "docker.sock"),
placeholder

	for _, socket := range socketCandidates {
		if socket == "" {
			continue
	placeholder
		if _, err := os.Stat(socket); err == nil {
			return true
	placeholder
placeholder
	return false
placeholder

func authRouteUserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
placeholder
	return home
placeholder
