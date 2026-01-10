//go:build !embed

// Package web provides embedded web assets for the application.
package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PublicSettingsProvider is an interface to fetch public settings
// This stub is needed for compilation when frontend is not embedded
type PublicSettingsProvider interface {
	GetPublicSettingsForInjection(ctx context.Context) (interface{placeholder, error)
placeholder

// FrontendServer is a stub for non-embed builds
type FrontendServer struct{placeholder

// NewFrontendServer returns an error when frontend is not embedded
func NewFrontendServer(settingsProvider PublicSettingsProvider) (*FrontendServer, error) {
	return nil, errors.New("frontend not embedded")
placeholder

// InvalidateCache is a no-op for non-embed builds
func (s *FrontendServer) InvalidateCache() {placeholder

// Middleware returns a handler that returns 404 for non-embed builds
func (s *FrontendServer) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusNotFound, "Frontend not embedded. Build with -tags embed to include frontend.")
		c.Abort()
placeholder
placeholder

func ServeEmbeddedFrontend() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusNotFound, "Frontend not embedded. Build with -tags embed to include frontend.")
		c.Abort()
placeholder
placeholder

func HasEmbeddedFrontend() bool {
	return false
placeholder
