package routes

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	servermiddleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type compositeRouteRepoStub struct {
	routes []service.CompositeModelRoute
placeholder

func (s compositeRouteRepoStub) ListByGroup(ctx context.Context, groupID int64, includeDisabled bool) ([]service.CompositeModelRoute, error) {
	routes := make([]service.CompositeModelRoute, 0, len(s.routes))
	for _, route := range s.routes {
		if route.GroupID != groupID {
			continue
	placeholder
		if !includeDisabled && !route.Enabled {
			continue
	placeholder
		routes = append(routes, route)
placeholder
	return routes, nil
placeholder

func (s compositeRouteRepoStub) Create(ctx context.Context, route *service.CompositeModelRoute) error {
	return nil
placeholder

func (s compositeRouteRepoStub) Update(ctx context.Context, route *service.CompositeModelRoute) error {
	return nil
placeholder

func (s compositeRouteRepoStub) Delete(ctx context.Context, id int64) error {
	return nil
placeholder

func (s compositeRouteRepoStub) DeleteByGroup(ctx context.Context, groupID int64) error {
	return nil
placeholder

func TestCompositeTargetPlatformMiddlewareResolvesModelAndRestoresBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.HandlerFunc(servermiddleware.APIKeyAuthMiddleware(func(c *gin.Context) {
		groupID := int64(1)
		c.Set(string(servermiddleware.ContextKeyAPIKey), &service.APIKey{
			GroupID: &groupID,
			Group:   &service.Group{Platform: service.PlatformCompositeplaceholder,
	placeholder)
		c.Next()
placeholder)))
	router.Use(compositeTargetPlatformMiddleware(nil))
	router.POST("/", func(c *gin.Context) {
		platform, ok := service.ResolvedTargetPlatformFromContext(c.Request.Context())
		require.True(t, ok)
		require.Equal(t, service.PlatformOpenAI, platform)

		body, err := io.ReadAll(c.Request.Body)
	placeholder
		require.JSONEq(t, `{"model":"gpt-5"placeholder`, string(body))
		c.Status(http.StatusNoContent)
placeholder)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"model":"gpt-5"placeholder`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
placeholder

func TestCompositeTargetPlatformMiddlewareUsesExplicitRouteAndRewritesBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	resolver := service.NewCompositeRouteResolver(compositeRouteRepoStub{
		routes: []service.CompositeModelRoute{
			{
				ID:             1,
				GroupID:        1,
				PublicModel:    "openrouter/gpt-5",
				MatchType:      service.CompositeRouteMatchExact,
				TargetPlatform: service.PlatformOpenAI,
				UpstreamModel:  "gpt-5",
				Endpoint:       service.CompositeRouteEndpointAny,
				Priority:       100,
				Enabled:        true,
		placeholder,
	placeholder,
placeholder)
	router.Use(gin.HandlerFunc(servermiddleware.APIKeyAuthMiddleware(func(c *gin.Context) {
		groupID := int64(1)
		c.Set(string(servermiddleware.ContextKeyAPIKey), &service.APIKey{
			GroupID: &groupID,
			Group:   &service.Group{ID: groupID, Platform: service.PlatformCompositeplaceholder,
	placeholder)
		c.Next()
placeholder)))
	router.Use(compositeTargetPlatformMiddleware(resolver))
	router.POST("/v1/chat/completions", func(c *gin.Context) {
		platform, ok := service.ResolvedTargetPlatformFromContext(c.Request.Context())
		require.True(t, ok)
		require.Equal(t, service.PlatformOpenAI, platform)

		upstreamModel, ok := service.ResolvedUpstreamModelFromContext(c.Request.Context())
		require.True(t, ok)
		require.Equal(t, "gpt-5", upstreamModel)

		body, err := io.ReadAll(c.Request.Body)
	placeholder
		require.JSONEq(t, `{"model":"gpt-5","messages":[]placeholder`, string(body))
		c.Status(http.StatusNoContent)
placeholder)

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", strings.NewReader(`{"model":"openrouter/gpt-5","messages":[]placeholder`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
placeholder

func TestCompositeTargetPlatformMiddlewareRewritesNestedLiveModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	resolver := service.NewCompositeRouteResolver(compositeRouteRepoStub{
		routes: []service.CompositeModelRoute{
			{
				ID:             1,
				GroupID:        1,
				PublicModel:    "live-alias",
				MatchType:      service.CompositeRouteMatchExact,
				TargetPlatform: service.PlatformOpenAI,
				UpstreamModel:  "gpt-live",
				Endpoint:       service.CompositeRouteEndpointAny,
				Priority:       100,
				Enabled:        true,
		placeholder,
	placeholder,
placeholder)
	router.Use(gin.HandlerFunc(servermiddleware.APIKeyAuthMiddleware(func(c *gin.Context) {
		groupID := int64(1)
		c.Set(string(servermiddleware.ContextKeyAPIKey), &service.APIKey{
			GroupID: &groupID,
			Group:   &service.Group{ID: groupID, Platform: service.PlatformCompositeplaceholder,
	placeholder)
		c.Next()
placeholder)))
	router.Use(compositeTargetPlatformMiddleware(resolver))
	router.POST("/backend-api/codex/realtime/calls", func(c *gin.Context) {
		platform, ok := service.ResolvedTargetPlatformFromContext(c.Request.Context())
		require.True(t, ok)
		require.Equal(t, service.PlatformOpenAI, platform)

		body, err := io.ReadAll(c.Request.Body)
	placeholder
		require.JSONEq(t, `{"session":{"model":"gpt-live"placeholder,"sdp":"v=0"placeholder`, string(body))
		c.Status(http.StatusNoContent)
placeholder)

	req := httptest.NewRequest(
		http.MethodPost,
		"/backend-api/codex/realtime/calls",
		strings.NewReader(`{"session":{"model":"live-alias"placeholder,"sdp":"v=0"placeholder`),
	)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
placeholder

func TestCompositeRequestModelFromMultipartLiveSession(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.WriteField("sdp", "v=0"))
	require.NoError(t, writer.WriteField("session", `{"model":"live-alias"placeholder`))
	require.NoError(t, writer.Close())

	require.Equal(t, "live-alias", compositeRequestModelFromBody(writer.FormDataContentType(), body.Bytes()))
placeholder

func TestCompositeCodexControlPathsUseResponsesRoutes(t *testing.T) {
	for _, path := range []string{
		"/v1/alpha/search",
		"/backend-api/codex/alpha/search",
		"/v1/live",
		"/backend-api/codex/realtime/calls",
placeholder {
		require.Equal(t, service.CompositeRouteEndpointResponses, compositeRouteEndpointForPath(path), "path=%s", path)
placeholder
placeholder

func TestCompositeTargetPlatformMiddlewareUsesExplicitRouteForMultipartImages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	resolver := service.NewCompositeRouteResolver(compositeRouteRepoStub{
		routes: []service.CompositeModelRoute{
			{
				ID:             1,
				GroupID:        1,
				PublicModel:    "image-alias",
				MatchType:      service.CompositeRouteMatchExact,
				TargetPlatform: service.PlatformOpenAI,
				UpstreamModel:  "gpt-image-1",
				Endpoint:       service.CompositeRouteEndpointImages,
				Priority:       100,
				Enabled:        true,
		placeholder,
	placeholder,
placeholder)
	router.Use(gin.HandlerFunc(servermiddleware.APIKeyAuthMiddleware(func(c *gin.Context) {
		groupID := int64(1)
		c.Set(string(servermiddleware.ContextKeyAPIKey), &service.APIKey{
			GroupID: &groupID,
			Group:   &service.Group{ID: groupID, Platform: service.PlatformCompositeplaceholder,
	placeholder)
		c.Next()
placeholder)))
	router.Use(compositeTargetPlatformMiddleware(resolver))
	router.POST("/v1/images/edits", func(c *gin.Context) {
		platform, ok := service.ResolvedTargetPlatformFromContext(c.Request.Context())
		require.True(t, ok)
		require.Equal(t, service.PlatformOpenAI, platform)

		upstreamModel, ok := service.ResolvedUpstreamModelFromContext(c.Request.Context())
		require.True(t, ok)
		require.Equal(t, "gpt-image-1", upstreamModel)

		publicModel, ok := service.RequestedPublicModelFromContext(c.Request.Context())
		require.True(t, ok)
		require.Equal(t, "image-alias", publicModel)

		body, err := io.ReadAll(c.Request.Body)
	placeholder
		require.Contains(t, string(body), "image-alias")
		c.Status(http.StatusNoContent)
placeholder)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.WriteField("model", "image-alias"))
	require.NoError(t, writer.WriteField("prompt", "draw"))
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/v1/images/edits", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
placeholder

func TestCompositeGeminiTargetPlatformMiddlewareUsesPathRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	resolver := service.NewCompositeRouteResolver(compositeRouteRepoStub{
		routes: []service.CompositeModelRoute{
			{
				ID:             1,
				GroupID:        1,
				PublicModel:    "openrouter/gemini-pro",
				MatchType:      service.CompositeRouteMatchExact,
				TargetPlatform: service.PlatformGemini,
				UpstreamModel:  "gemini-2.5-pro",
				Endpoint:       service.CompositeRouteEndpointGemini,
				Priority:       100,
				Enabled:        true,
		placeholder,
	placeholder,
placeholder)
	router.Use(gin.HandlerFunc(servermiddleware.APIKeyAuthMiddleware(func(c *gin.Context) {
		groupID := int64(1)
		c.Set(string(servermiddleware.ContextKeyAPIKey), &service.APIKey{
			GroupID: &groupID,
			Group:   &service.Group{ID: groupID, Platform: service.PlatformCompositeplaceholder,
	placeholder)
		c.Next()
placeholder)))
	router.Use(compositeGeminiTargetPlatformMiddleware(resolver))
	router.POST("/v1beta/models/*modelAction", func(c *gin.Context) {
		platform, ok := service.ResolvedTargetPlatformFromContext(c.Request.Context())
		require.True(t, ok)
		require.Equal(t, service.PlatformGemini, platform)

		upstreamModel, ok := service.ResolvedUpstreamModelFromContext(c.Request.Context())
		require.True(t, ok)
		require.Equal(t, "gemini-2.5-pro", upstreamModel)
		c.Status(http.StatusNoContent)
placeholder)

	req := httptest.NewRequest(http.MethodPost, "/v1beta/models/openrouter/gemini-pro:generateContent", strings.NewReader(`{"contents":[]placeholder`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
placeholder
