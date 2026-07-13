package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/servertiming"
	"github.com/gin-gonic/gin"
)

func runServerTimingRequest(
	t *testing.T,
	enabled bool,
	path string,
	marker string,
	role string,
	handler gin.HandlerFunc,
) *httptest.ResponseRecorder {
placeholder
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(ServerTiming(enabled))
	engine.Any("/*path", func(c *gin.Context) {
		if role != "" {
			c.Set(string(ContextKeyUserRole), role)
	placeholder
		handler(c)
placeholder)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, path, nil)
	if marker != "" {
		request.Header.Set(servertiming.AdminUIHeader, marker)
placeholder
	engine.ServeHTTP(recorder, request)
	return recorder
placeholder

func TestServerTimingScopesAndRoleGate(t *testing.T) {
	tests := []struct {
		name       string
		enabled    bool
		path       string
		marker     string
		role       string
		wantHeader bool
placeholder{
		{name: "disabled", enabled: false, path: "/api/v1/admin/users", role: "admin"placeholder,
		{name: "admin API path", enabled: true, path: "/api/v1/admin/users", role: "admin", wantHeader: trueplaceholder,
		{name: "shared API marked by admin UI", enabled: true, path: "/api/v1/groups/available", marker: "1", role: "admin", wantHeader: trueplaceholder,
		{name: "non admin role", enabled: true, path: "/api/v1/groups/available", marker: "1", role: "user"placeholder,
		{name: "unauthenticated public request", enabled: true, path: "/api/v1/settings/public", marker: "1"placeholder,
		{name: "unmarked shared API", enabled: true, path: "/api/v1/groups/available", role: "admin"placeholder,
		{name: "invalid marker", enabled: true, path: "/api/v1/groups/available", marker: "true", role: "admin"placeholder,
		{name: "admin prefix boundary", enabled: true, path: "/api/v1/administrator", role: "admin"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := runServerTimingRequest(t, tt.enabled, tt.path, tt.marker, tt.role, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": trueplaceholder)
		placeholder)
			header := recorder.Header().Get(servertiming.HeaderName)
			if tt.wantHeader && header == "" {
				t.Fatalf("%s header missing", servertiming.HeaderName)
		placeholder
			if !tt.wantHeader && header != "" {
				t.Fatalf("unexpected %s header: %q", servertiming.HeaderName, header)
		placeholder
			if header != "" && (!strings.Contains(header, "total;dur=") || !strings.Contains(header, `cache;desc="bypass"`)) {
				t.Fatalf("incomplete timing header: %q", header)
		placeholder
	placeholder)
placeholder
placeholder

func TestServerTimingCollectorIsRequestScoped(t *testing.T) {
	active := false
	recorder := runServerTimingRequest(t, true, "/api/v1/keys", "1", "admin", func(c *gin.Context) {
		active = servertiming.Active(c.Request.Context())
		c.Status(http.StatusNoContent)
placeholder)
	if !active {
		t.Fatal("collector was not attached to marked request context")
placeholder
	if recorder.Header().Get(servertiming.HeaderName) == "" {
		t.Fatal("timing header missing from status-only response")
placeholder
placeholder

func TestServerTimingFinalizesBeforeEarlyCommit(t *testing.T) {
	recorder := runServerTimingRequest(t, true, "/api/v1/admin/stream", "", "admin", func(c *gin.Context) {
		c.Status(http.StatusAccepted)
		c.Writer.WriteHeaderNow()
placeholder)
	if got := recorder.Header().Get(servertiming.HeaderName); got == "" {
		t.Fatal("timing header was not written before response commit")
placeholder
placeholder

func TestServerTimingFinalizesOnFlush(t *testing.T) {
	recorder := runServerTimingRequest(t, true, "/api/v1/admin/export", "", "admin", func(c *gin.Context) {
		c.Writer.Flush()
placeholder)
	if got := recorder.Header().Get(servertiming.HeaderName); got == "" {
		t.Fatal("timing header was not written before stream flush")
placeholder
placeholder

func TestServerTimingStatusResponses(t *testing.T) {
	tests := []struct {
		name   string
		status int
placeholder{
		{name: "not modified", status: http.StatusNotModifiedplaceholder,
		{name: "internal error", status: http.StatusInternalServerErrorplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := runServerTimingRequest(t, true, "/api/v1/admin/test", "", "admin", func(c *gin.Context) {
				c.Status(tt.status)
		placeholder)
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.status)
		placeholder
			if got := recorder.Header().Get(servertiming.HeaderName); got == "" {
				t.Fatalf("timing header missing from status %d response", tt.status)
		placeholder
	placeholder)
placeholder
placeholder

func TestServerTimingResponseWriterUnwraps(t *testing.T) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	baseWriter := c.Writer
	writer := &serverTimingResponseWriter{ResponseWriter: baseWriterplaceholder
	if got := writer.Unwrap(); got != baseWriter {
		t.Fatalf("Unwrap() = %T, want original Gin writer", got)
placeholder
placeholder

func TestServerTimingCacheOutcome(t *testing.T) {
	tests := []struct {
		name       string
		headerName string
		value      string
		want       string
placeholder{
		{name: "snapshot hit", headerName: snapshotCacheHeader, value: "hit", want: "hit"placeholder,
		{name: "usage miss", headerName: usageCacheHeader, value: "MISS", want: "miss"placeholder,
		{name: "invalid", headerName: snapshotCacheHeader, value: "stale", want: "bypass"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := runServerTimingRequest(t, true, "/api/v1/admin/dashboard", "", "admin", func(c *gin.Context) {
				c.Header(tt.headerName, tt.value)
				c.JSON(http.StatusOK, gin.H{"ok": trueplaceholder)
		placeholder)
			want := `cache;desc="` + tt.want + `"`
			if got := recorder.Header().Get(servertiming.HeaderName); !strings.Contains(got, want) {
				t.Fatalf("timing header %q does not contain %q", got, want)
		placeholder
	placeholder)
placeholder
placeholder

func TestServerTimingResponseHeaderForWebSocket(t *testing.T) {
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/ops/ws/qps", nil)
	collector := servertiming.New(time.Now())
	c.Request = c.Request.WithContext(servertiming.WithCollector(c.Request.Context(), collector))
	c.Set(string(ContextKeyUserRole), "admin")

	header := ServerTimingResponseHeader(c)
	if header.Get(servertiming.HeaderName) == "" {
		t.Fatal("WebSocket response header missing timing value")
placeholder

	c.Set(string(ContextKeyUserRole), "user")
	if got := ServerTimingResponseHeader(c); got != nil {
		t.Fatalf("non-admin WebSocket received timing header: %#v", got)
placeholder
placeholder
