//go:build embed

package web

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed all:dist
var frontendFS embed.FS

// PublicSettingsProvider is an interface to fetch public settings
type PublicSettingsProvider interface {
	GetPublicSettingsForInjection(ctx context.Context) (any, error)
placeholder

// FrontendServer serves the embedded frontend with settings injection
type FrontendServer struct {
	distFS     fs.FS
	fileServer http.Handler
	baseHTML   []byte
	cache      *HTMLCache
	settings   PublicSettingsProvider
placeholder

// NewFrontendServer creates a new frontend server with settings injection
func NewFrontendServer(settingsProvider PublicSettingsProvider) (*FrontendServer, error) {
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		return nil, err
placeholder

	// Read base HTML once
	file, err := distFS.Open("index.html")
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = file.Close() placeholder()

	baseHTML, err := io.ReadAll(file)
	if err != nil {
		return nil, err
placeholder

	cache := NewHTMLCache()
	cache.SetBaseHTML(baseHTML)

	return &FrontendServer{
		distFS:     distFS,
		fileServer: http.FileServer(http.FS(distFS)),
		baseHTML:   baseHTML,
		cache:      cache,
		settings:   settingsProvider,
placeholder, nil
placeholder

// InvalidateCache invalidates the HTML cache (call when settings change)
func (s *FrontendServer) InvalidateCache() {
	if s != nil && s.cache != nil {
		s.cache.Invalidate()
placeholder
placeholder

// Middleware returns the Gin middleware handler
func (s *FrontendServer) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip API routes
		if strings.HasPrefix(path, "/api/") ||
			strings.HasPrefix(path, "/v1/") ||
			strings.HasPrefix(path, "/v1beta/") ||
			strings.HasPrefix(path, "/antigravity/") ||
			strings.HasPrefix(path, "/setup/") ||
			path == "/health" ||
			path == "/responses" {
			c.Next()
			return
	placeholder

		cleanPath := strings.TrimPrefix(path, "/")
		if cleanPath == "" {
			cleanPath = "index.html"
	placeholder

		// For index.html or SPA routes, serve with injected settings
		if cleanPath == "index.html" || !s.fileExists(cleanPath) {
			s.serveIndexHTML(c)
			return
	placeholder

		// Serve static files normally
		s.fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
placeholder
placeholder

func (s *FrontendServer) fileExists(path string) bool {
	file, err := s.distFS.Open(path)
	if err != nil {
		return false
placeholder
	_ = file.Close()
	return true
placeholder

func (s *FrontendServer) serveIndexHTML(c *gin.Context) {
	// Check cache first
	cached := s.cache.Get()
	if cached != nil {
		// Check If-None-Match for 304 response
		if match := c.GetHeader("If-None-Match"); match == cached.ETag {
			c.Status(http.StatusNotModified)
			c.Abort()
			return
	placeholder

		c.Header("ETag", cached.ETag)
		c.Header("Cache-Control", "no-cache") // Must revalidate
		c.Data(http.StatusOK, "text/html; charset=utf-8", cached.Content)
		c.Abort()
		return
placeholder

	// Cache miss - fetch settings and render
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	settings, err := s.settings.GetPublicSettingsForInjection(ctx)
	if err != nil {
		// Fallback: serve without injection
		c.Data(http.StatusOK, "text/html; charset=utf-8", s.baseHTML)
		c.Abort()
		return
placeholder

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		// Fallback: serve without injection
		c.Data(http.StatusOK, "text/html; charset=utf-8", s.baseHTML)
		c.Abort()
		return
placeholder

	rendered := s.injectSettings(settingsJSON)
	s.cache.Set(rendered, settingsJSON)

	cached = s.cache.Get()
	if cached != nil {
		c.Header("ETag", cached.ETag)
placeholder
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/html; charset=utf-8", rendered)
	c.Abort()
placeholder

func (s *FrontendServer) injectSettings(settingsJSON []byte) []byte {
	// Create the script tag to inject
	script := []byte(`<script>window.__APP_CONFIG__=` + string(settingsJSON) + `;</script>`)

	// Inject before </head>
	headClose := []byte("</head>")
	return bytes.Replace(s.baseHTML, headClose, append(script, headClose...), 1)
placeholder

// ServeEmbeddedFrontend returns a middleware for serving embedded frontend
// This is the legacy function for backward compatibility when no settings provider is available
func ServeEmbeddedFrontend() gin.HandlerFunc {
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		panic("failed to get dist subdirectory: " + err.Error())
placeholder
	fileServer := http.FileServer(http.FS(distFS))

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.HasPrefix(path, "/api/") ||
			strings.HasPrefix(path, "/v1/") ||
			strings.HasPrefix(path, "/v1beta/") ||
			strings.HasPrefix(path, "/antigravity/") ||
			strings.HasPrefix(path, "/setup/") ||
			path == "/health" ||
			path == "/responses" {
			c.Next()
			return
	placeholder

		cleanPath := strings.TrimPrefix(path, "/")
		if cleanPath == "" {
			cleanPath = "index.html"
	placeholder

		if file, err := distFS.Open(cleanPath); err == nil {
			_ = file.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
	placeholder

		serveIndexHTML(c, distFS)
placeholder
placeholder

func serveIndexHTML(c *gin.Context, fsys fs.FS) {
	file, err := fsys.Open("index.html")
	if err != nil {
		c.String(http.StatusNotFound, "Frontend not found")
		c.Abort()
		return
placeholder
	defer func() { _ = file.Close() placeholder()

	content, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read index.html")
		c.Abort()
		return
placeholder

	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	c.Abort()
placeholder

func HasEmbeddedFrontend() bool {
	_, err := frontendFS.ReadFile("dist/index.html")
	return err == nil
placeholder
