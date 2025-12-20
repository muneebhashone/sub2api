package setup

import (
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"sync"
	"time"

	"sub2api/internal/pkg/response"
	"sub2api/internal/pkg/sysutil"

	"github.com/gin-gonic/gin"
)

// installMutex prevents concurrent installation attempts (TOCTOU protection)
var installMutex sync.Mutex

// RegisterRoutes registers setup wizard routes
func RegisterRoutes(r *gin.Engine) {
	setup := r.Group("/setup")
	{
		// Status endpoint is always accessible (read-only)
		setup.GET("/status", getStatus)

		// All modification endpoints are protected by setupGuard
		protected := setup.Group("")
		protected.Use(setupGuard())
		{
			protected.POST("/test-db", testDatabase)
			protected.POST("/test-redis", testRedis)
			protected.POST("/install", install)
	placeholder
placeholder
placeholder

// SetupStatus represents the current setup state
type SetupStatus struct {
	NeedsSetup bool   `json:"needs_setup"`
	Step       string `json:"step"`
placeholder

// getStatus returns the current setup status
func getStatus(c *gin.Context) {
	response.Success(c, SetupStatus{
		NeedsSetup: NeedsSetup(),
		Step:       "welcome",
placeholder)
placeholder

// setupGuard middleware ensures setup endpoints are only accessible during setup mode
func setupGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !NeedsSetup() {
			response.Error(c, http.StatusForbidden, "Setup is not allowed: system is already installed")
			c.Abort()
			return
	placeholder
		c.Next()
placeholder
placeholder

// validateHostname checks if a hostname/IP is safe (no injection characters)
func validateHostname(host string) bool {
	// Allow only alphanumeric, dots, hyphens, and colons (for IPv6)
	validHost := regexp.MustCompile(`^[a-zA-Z0-9.\-:]+$`)
	return validHost.MatchString(host) && len(host) <= 253
placeholder

// validateDBName checks if database name is safe
func validateDBName(name string) bool {
	// Allow only alphanumeric and underscores, starting with letter
	validName := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	return validName.MatchString(name) && len(name) <= 63
placeholder

// validateUsername checks if username is safe
func validateUsername(name string) bool {
	// Allow only alphanumeric and underscores
	validName := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validName.MatchString(name) && len(name) <= 63
placeholder

// validateEmail checks if email format is valid
func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil && len(email) <= 254
placeholder

// validatePassword checks password strength
func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
placeholder
	if len(password) > 128 {
		return fmt.Errorf("password must be at most 128 characters")
placeholder
	return nil
placeholder

// validatePort checks if port is in valid range
func validatePort(port int) bool {
	return port > 0 && port <= 65535
placeholder

// validateSSLMode checks if SSL mode is valid
func validateSSLMode(mode string) bool {
	validModes := map[string]bool{
		"disable": true, "require": true, "verify-ca": true, "verify-full": true,
placeholder
	return validModes[mode]
placeholder

// TestDatabaseRequest represents database test request
type TestDatabaseRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password"`
	DBName   string `json:"dbname" binding:"required"`
	SSLMode  string `json:"sslmode"`
placeholder

// testDatabase tests database connection
func testDatabase(c *gin.Context) {
	var req TestDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
placeholder

	// Security: Validate all inputs to prevent injection attacks
	if !validateHostname(req.Host) {
		response.Error(c, http.StatusBadRequest, "Invalid hostname format")
		return
placeholder
	if !validatePort(req.Port) {
		response.Error(c, http.StatusBadRequest, "Invalid port number")
		return
placeholder
	if !validateUsername(req.User) {
		response.Error(c, http.StatusBadRequest, "Invalid username format")
		return
placeholder
	if !validateDBName(req.DBName) {
		response.Error(c, http.StatusBadRequest, "Invalid database name format")
		return
placeholder

	if req.SSLMode == "" {
		req.SSLMode = "disable"
placeholder
	if !validateSSLMode(req.SSLMode) {
		response.Error(c, http.StatusBadRequest, "Invalid SSL mode")
		return
placeholder

	cfg := &DatabaseConfig{
		Host:     req.Host,
		Port:     req.Port,
		User:     req.User,
		Password: req.Password,
		DBName:   req.DBName,
		SSLMode:  req.SSLMode,
placeholder

	if err := TestDatabaseConnection(cfg); err != nil {
		response.Error(c, http.StatusBadRequest, "Connection failed: "+err.Error())
		return
placeholder

	response.Success(c, gin.H{"message": "Connection successful"placeholder)
placeholder

// TestRedisRequest represents Redis test request
type TestRedisRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Password string `json:"password"`
	DB       int    `json:"db"`
placeholder

// testRedis tests Redis connection
func testRedis(c *gin.Context) {
	var req TestRedisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
placeholder

	// Security: Validate inputs
	if !validateHostname(req.Host) {
		response.Error(c, http.StatusBadRequest, "Invalid hostname format")
		return
placeholder
	if !validatePort(req.Port) {
		response.Error(c, http.StatusBadRequest, "Invalid port number")
		return
placeholder
	if req.DB < 0 || req.DB > 15 {
		response.Error(c, http.StatusBadRequest, "Invalid Redis database number (0-15)")
		return
placeholder

	cfg := &RedisConfig{
		Host:     req.Host,
		Port:     req.Port,
		Password: req.Password,
		DB:       req.DB,
placeholder

	if err := TestRedisConnection(cfg); err != nil {
		response.Error(c, http.StatusBadRequest, "Connection failed: "+err.Error())
		return
placeholder

	response.Success(c, gin.H{"message": "Connection successful"placeholder)
placeholder

// InstallRequest represents installation request
type InstallRequest struct {
	Database DatabaseConfig `json:"database" binding:"required"`
	Redis    RedisConfig    `json:"redis" binding:"required"`
	Admin    AdminConfig    `json:"admin" binding:"required"`
	Server   ServerConfig   `json:"server"`
placeholder

// install performs the installation
func install(c *gin.Context) {
	// TOCTOU Protection: Acquire mutex to prevent concurrent installation
	installMutex.Lock()
	defer installMutex.Unlock()

	// Double-check after acquiring lock
	if !NeedsSetup() {
		response.Error(c, http.StatusForbidden, "Setup is not allowed: system is already installed")
		return
placeholder

	var req InstallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
placeholder

	// ========== COMPREHENSIVE INPUT VALIDATION ==========
	// Database validation
	if !validateHostname(req.Database.Host) {
		response.Error(c, http.StatusBadRequest, "Invalid database hostname")
		return
placeholder
	if !validatePort(req.Database.Port) {
		response.Error(c, http.StatusBadRequest, "Invalid database port")
		return
placeholder
	if !validateUsername(req.Database.User) {
		response.Error(c, http.StatusBadRequest, "Invalid database username")
		return
placeholder
	if !validateDBName(req.Database.DBName) {
		response.Error(c, http.StatusBadRequest, "Invalid database name")
		return
placeholder

	// Redis validation
	if !validateHostname(req.Redis.Host) {
		response.Error(c, http.StatusBadRequest, "Invalid Redis hostname")
		return
placeholder
	if !validatePort(req.Redis.Port) {
		response.Error(c, http.StatusBadRequest, "Invalid Redis port")
		return
placeholder
	if req.Redis.DB < 0 || req.Redis.DB > 15 {
		response.Error(c, http.StatusBadRequest, "Invalid Redis database number")
		return
placeholder

	// Admin validation
	if !validateEmail(req.Admin.Email) {
		response.Error(c, http.StatusBadRequest, "Invalid admin email format")
		return
placeholder
	if err := validatePassword(req.Admin.Password); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
placeholder

	// Server validation
	if req.Server.Port != 0 && !validatePort(req.Server.Port) {
		response.Error(c, http.StatusBadRequest, "Invalid server port")
		return
placeholder

	// ========== SET DEFAULTS ==========
	if req.Database.SSLMode == "" {
		req.Database.SSLMode = "disable"
placeholder
	if !validateSSLMode(req.Database.SSLMode) {
		response.Error(c, http.StatusBadRequest, "Invalid SSL mode")
		return
placeholder
	if req.Server.Host == "" {
		req.Server.Host = "0.0.0.0"
placeholder
	if req.Server.Port == 0 {
		req.Server.Port = 8080
placeholder
	if req.Server.Mode == "" {
		req.Server.Mode = "release"
placeholder
	// Validate server mode
	if req.Server.Mode != "release" && req.Server.Mode != "debug" {
		response.Error(c, http.StatusBadRequest, "Invalid server mode (must be 'release' or 'debug')")
		return
placeholder

	// Trim whitespace from string inputs
	req.Admin.Email = strings.TrimSpace(req.Admin.Email)
	req.Database.Host = strings.TrimSpace(req.Database.Host)
	req.Database.User = strings.TrimSpace(req.Database.User)
	req.Database.DBName = strings.TrimSpace(req.Database.DBName)
	req.Redis.Host = strings.TrimSpace(req.Redis.Host)

	cfg := &SetupConfig{
		Database: req.Database,
		Redis:    req.Redis,
		Admin:    req.Admin,
		Server:   req.Server,
		JWT: JWTConfig{
			ExpireHour: 24,
	placeholder,
placeholder

	if err := Install(cfg); err != nil {
		response.Error(c, http.StatusInternalServerError, "Installation failed: "+err.Error())
		return
placeholder

	// Schedule service restart in background after sending response
	// This ensures the client receives the success response before the service restarts
	go func() {
		// Wait a moment to ensure the response is sent
		time.Sleep(500 * time.Millisecond)
		sysutil.RestartServiceAsync()
placeholder()

	response.Success(c, gin.H{
		"message": "Installation completed successfully. Service will restart automatically.",
		"restart": true,
placeholder)
placeholder
