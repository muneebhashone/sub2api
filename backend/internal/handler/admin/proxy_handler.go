package admin

import (
	"strconv"

	"sub2api/internal/pkg/response"
	"sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ProxyHandler handles admin proxy management
type ProxyHandler struct {
	adminService service.AdminService
placeholder

// NewProxyHandler creates a new admin proxy handler
func NewProxyHandler(adminService service.AdminService) *ProxyHandler {
	return &ProxyHandler{
		adminService: adminService,
placeholder
placeholder

// CreateProxyRequest represents create proxy request
type CreateProxyRequest struct {
	Name     string `json:"name" binding:"required"`
	Protocol string `json:"protocol" binding:"required,oneof=http https socks5"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required,min=1,max=65535"`
	Username string `json:"username"`
	Password string `json:"password"`
placeholder

// UpdateProxyRequest represents update proxy request
type UpdateProxyRequest struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol" binding:"omitempty,oneof=http https socks5"`
	Host     string `json:"host"`
	Port     int    `json:"port" binding:"omitempty,min=1,max=65535"`
	Username string `json:"username"`
	Password string `json:"password"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive"`
placeholder

// List handles listing all proxies with pagination
// GET /api/v1/admin/proxies
func (h *ProxyHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	protocol := c.Query("protocol")
	status := c.Query("status")
	search := c.Query("search")

	proxies, total, err := h.adminService.ListProxies(c.Request.Context(), page, pageSize, protocol, status, search)
	if err != nil {
		response.InternalError(c, "Failed to list proxies: "+err.Error())
		return
placeholder

	response.Paginated(c, proxies, total, page, pageSize)
placeholder

// GetAll handles getting all active proxies without pagination
// GET /api/v1/admin/proxies/all
// Optional query param: with_count=true to include account count per proxy
func (h *ProxyHandler) GetAll(c *gin.Context) {
	withCount := c.Query("with_count") == "true"

	if withCount {
		proxies, err := h.adminService.GetAllProxiesWithAccountCount(c.Request.Context())
		if err != nil {
			response.InternalError(c, "Failed to get proxies: "+err.Error())
			return
	placeholder
		response.Success(c, proxies)
		return
placeholder

	proxies, err := h.adminService.GetAllProxies(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Failed to get proxies: "+err.Error())
		return
placeholder

	response.Success(c, proxies)
placeholder

// GetByID handles getting a proxy by ID
// GET /api/v1/admin/proxies/:id
func (h *ProxyHandler) GetByID(c *gin.Context) {
	proxyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid proxy ID")
		return
placeholder

	proxy, err := h.adminService.GetProxy(c.Request.Context(), proxyID)
	if err != nil {
		response.NotFound(c, "Proxy not found")
		return
placeholder

	response.Success(c, proxy)
placeholder

// Create handles creating a new proxy
// POST /api/v1/admin/proxies
func (h *ProxyHandler) Create(c *gin.Context) {
	var req CreateProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	proxy, err := h.adminService.CreateProxy(c.Request.Context(), &service.CreateProxyInput{
		Name:     req.Name,
		Protocol: req.Protocol,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
placeholder)
	if err != nil {
		response.BadRequest(c, "Failed to create proxy: "+err.Error())
		return
placeholder

	response.Success(c, proxy)
placeholder

// Update handles updating a proxy
// PUT /api/v1/admin/proxies/:id
func (h *ProxyHandler) Update(c *gin.Context) {
	proxyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid proxy ID")
		return
placeholder

	var req UpdateProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	proxy, err := h.adminService.UpdateProxy(c.Request.Context(), proxyID, &service.UpdateProxyInput{
		Name:     req.Name,
		Protocol: req.Protocol,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Status:   req.Status,
placeholder)
	if err != nil {
		response.InternalError(c, "Failed to update proxy: "+err.Error())
		return
placeholder

	response.Success(c, proxy)
placeholder

// Delete handles deleting a proxy
// DELETE /api/v1/admin/proxies/:id
func (h *ProxyHandler) Delete(c *gin.Context) {
	proxyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid proxy ID")
		return
placeholder

	err = h.adminService.DeleteProxy(c.Request.Context(), proxyID)
	if err != nil {
		response.InternalError(c, "Failed to delete proxy: "+err.Error())
		return
placeholder

	response.Success(c, gin.H{"message": "Proxy deleted successfully"placeholder)
placeholder

// Test handles testing proxy connectivity
// POST /api/v1/admin/proxies/:id/test
func (h *ProxyHandler) Test(c *gin.Context) {
	proxyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid proxy ID")
		return
placeholder

	result, err := h.adminService.TestProxy(c.Request.Context(), proxyID)
	if err != nil {
		response.InternalError(c, "Failed to test proxy: "+err.Error())
		return
placeholder

	response.Success(c, result)
placeholder

// GetStats handles getting proxy statistics
// GET /api/v1/admin/proxies/:id/stats
func (h *ProxyHandler) GetStats(c *gin.Context) {
	proxyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid proxy ID")
		return
placeholder

	// Return mock data for now
	_ = proxyID
	response.Success(c, gin.H{
		"total_accounts":  0,
		"active_accounts": 0,
		"total_requests":  0,
		"success_rate":    100.0,
		"average_latency": 0,
placeholder)
placeholder

// GetProxyAccounts handles getting accounts using a proxy
// GET /api/v1/admin/proxies/:id/accounts
func (h *ProxyHandler) GetProxyAccounts(c *gin.Context) {
	proxyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid proxy ID")
		return
placeholder

	page, pageSize := response.ParsePagination(c)

	accounts, total, err := h.adminService.GetProxyAccounts(c.Request.Context(), proxyID, page, pageSize)
	if err != nil {
		response.InternalError(c, "Failed to get proxy accounts: "+err.Error())
		return
placeholder

	response.Paginated(c, accounts, total, page, pageSize)
placeholder


// BatchCreateProxyItem represents a single proxy in batch create request
type BatchCreateProxyItem struct {
	Protocol string `json:"protocol" binding:"required,oneof=http https socks5"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required,min=1,max=65535"`
	Username string `json:"username"`
	Password string `json:"password"`
placeholder

// BatchCreateRequest represents batch create proxies request
type BatchCreateRequest struct {
	Proxies []BatchCreateProxyItem `json:"proxies" binding:"required,min=1"`
placeholder

// BatchCreate handles batch creating proxies
// POST /api/v1/admin/proxies/batch
func (h *ProxyHandler) BatchCreate(c *gin.Context) {
	var req BatchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	created := 0
	skipped := 0

	for _, item := range req.Proxies {
		// Check for duplicates (same host, port, username, password)
		exists, err := h.adminService.CheckProxyExists(c.Request.Context(), item.Host, item.Port, item.Username, item.Password)
		if err != nil {
			response.InternalError(c, "Failed to check proxy existence: "+err.Error())
			return
	placeholder

		if exists {
			skipped++
			continue
	placeholder

		// Create proxy with default name
		_, err = h.adminService.CreateProxy(c.Request.Context(), &service.CreateProxyInput{
			Name:     "default",
			Protocol: item.Protocol,
			Host:     item.Host,
			Port:     item.Port,
			Username: item.Username,
			Password: item.Password,
	placeholder)
		if err != nil {
			// If creation fails due to duplicate, count as skipped
			skipped++
			continue
	placeholder

		created++
placeholder

	response.Success(c, gin.H{
		"created": created,
		"skipped": skipped,
placeholder)
placeholder
