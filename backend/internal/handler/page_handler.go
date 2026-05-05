package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

var validSlugPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

const maxPageFileSize = 1 << 20 // 1MB

type PageHandler struct {
	pagesDir string
placeholder

func NewPageHandler(dataDir string) *PageHandler {
	pagesDir := filepath.Join(dataDir, "pages")
	_ = os.MkdirAll(pagesDir, 0755)
	return &PageHandler{pagesDir: pagesDirplaceholder
placeholder

// GetPageContent serves raw markdown content for a given slug.
// GET /api/v1/pages/:slug
func (h *PageHandler) GetPageContent(c *gin.Context) {
	slug := c.Param("slug")
	if !validSlugPattern.MatchString(slug) || len(slug) > 64 {
		response.BadRequest(c, "Invalid page slug")
		return
placeholder

	filePath := filepath.Join(h.pagesDir, slug+".md")
	cleaned := filepath.Clean(filePath)
	if !strings.HasPrefix(cleaned, filepath.Clean(h.pagesDir)) {
		response.BadRequest(c, "Invalid page slug")
		return
placeholder

	info, err := os.Stat(cleaned)
	if err != nil || info.IsDir() {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"placeholder)
		return
placeholder
	if info.Size() > maxPageFileSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "page too large"placeholder)
		return
placeholder

	content, err := os.ReadFile(cleaned)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read page"placeholder)
		return
placeholder

	c.Data(http.StatusOK, "text/markdown; charset=utf-8", content)
placeholder

// ListPages returns available page slugs.
// GET /api/v1/pages
func (h *PageHandler) ListPages(c *gin.Context) {
	entries, err := os.ReadDir(h.pagesDir)
	if err != nil {
		response.Success(c, []string{placeholder)
		return
placeholder

	slugs := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
	placeholder
		name := e.Name()
		if strings.HasSuffix(name, ".md") {
			slugs = append(slugs, strings.TrimSuffix(name, ".md"))
	placeholder
placeholder
	response.Success(c, slugs)
placeholder

// ServePageImage serves images from data/pages/{slugplaceholder/ directory.
// GET /api/v1/pages/:slug/images/*filename
func (h *PageHandler) ServePageImage(c *gin.Context) {
	slug := c.Param("slug")
	filename := c.Param("filename")
	filename = strings.TrimPrefix(filename, "/")

	if !validSlugPattern.MatchString(slug) || len(slug) > 64 {
		c.Status(http.StatusNotFound)
		return
placeholder
	if filename == "" || strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		c.Status(http.StatusNotFound)
		return
placeholder

	imagesDir := filepath.Join(h.pagesDir, slug)
	filePath := filepath.Join(imagesDir, filename)
	cleaned := filepath.Clean(filePath)
	if !strings.HasPrefix(cleaned, filepath.Clean(imagesDir)) {
		c.Status(http.StatusNotFound)
		return
placeholder

	info, err := os.Stat(cleaned)
	if err != nil || info.IsDir() {
		c.Status(http.StatusNotFound)
		return
placeholder

	c.File(cleaned)
placeholder

// RegisterPageRoutes registers page routes on a router group.
func RegisterPageRoutes(v1 *gin.RouterGroup, dataDir string) {
	h := NewPageHandler(dataDir)
	pages := v1.Group("/pages")
	{
		pages.GET("", h.ListPages)
		pages.GET("/:slug", h.GetPageContent)
		pages.GET("/:slug/images/*filename", h.ServePageImage)
placeholder
placeholder
