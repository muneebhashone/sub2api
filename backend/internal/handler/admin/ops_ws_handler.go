package admin

import (
	"context"
	"encoding/json"
	"math"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type OpsWSProxyConfig struct {
	TrustProxy     bool
	TrustedProxies []netip.Prefix
	OriginPolicy   string
placeholder

const (
	envOpsWSTrustProxy     = "OPS_WS_TRUST_PROXY"
	envOpsWSTrustedProxies = "OPS_WS_TRUSTED_PROXIES"
	envOpsWSOriginPolicy   = "OPS_WS_ORIGIN_POLICY"
	envOpsWSMaxConns       = "OPS_WS_MAX_CONNS"
	envOpsWSMaxConnsPerIP  = "OPS_WS_MAX_CONNS_PER_IP"
)

const (
	OriginPolicyStrict     = "strict"
	OriginPolicyPermissive = "permissive"
)

var opsWSProxyConfig = loadOpsWSProxyConfigFromEnv()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return isAllowedOpsWSOrigin(r)
placeholder,
	// Subprotocol negotiation:
	// - The frontend passes ["sub2api-admin", "jwt.<token>"].
	// - We always select "sub2api-admin" so the token is never echoed back in the handshake response.
	Subprotocols: []string{"sub2api-admin"placeholder,
placeholder

const (
	qpsWSPushInterval       = 2 * time.Second
	qpsWSRefreshInterval    = 5 * time.Second
	qpsWSRequestCountWindow = 1 * time.Minute

	defaultMaxWSConns      = 100
	defaultMaxWSConnsPerIP = 20
)

var wsConnCount atomic.Int32
var wsConnCountByIP sync.Map // map[string]*atomic.Int32

const qpsWSIdleStopDelay = 30 * time.Second

const (
	opsWSCloseRealtimeDisabled = 4001
)

var qpsWSIdleStopMu sync.Mutex
var qpsWSIdleStopTimer *time.Timer

func cancelQPSWSIdleStop() {
	qpsWSIdleStopMu.Lock()
	if qpsWSIdleStopTimer != nil {
		qpsWSIdleStopTimer.Stop()
		qpsWSIdleStopTimer = nil
placeholder
	qpsWSIdleStopMu.Unlock()
placeholder

func scheduleQPSWSIdleStop() {
	qpsWSIdleStopMu.Lock()
	if qpsWSIdleStopTimer != nil {
		qpsWSIdleStopMu.Unlock()
		return
placeholder
	qpsWSIdleStopTimer = time.AfterFunc(qpsWSIdleStopDelay, func() {
		// Only stop if truly idle at fire time.
		if wsConnCount.Load() == 0 {
			qpsWSCache.Stop()
	placeholder
		qpsWSIdleStopMu.Lock()
		qpsWSIdleStopTimer = nil
		qpsWSIdleStopMu.Unlock()
placeholder)
	qpsWSIdleStopMu.Unlock()
placeholder

type opsWSRuntimeLimits struct {
	MaxConns      int32
	MaxConnsPerIP int32
placeholder

var opsWSLimits = loadOpsWSRuntimeLimitsFromEnv()

const (
	qpsWSWriteTimeout = 10 * time.Second
	qpsWSPongWait     = 60 * time.Second
	qpsWSPingInterval = 30 * time.Second

	// We don't expect clients to send application messages; we only read to process control frames (Pong/Close).
	qpsWSMaxReadBytes = 1024
)

type opsWSQPSCache struct {
	refreshInterval    time.Duration
	requestCountWindow time.Duration

	lastUpdatedUnixNano atomic.Int64
	payload             atomic.Value // []byte

	opsService *service.OpsService
	cancel     context.CancelFunc
	done       chan struct{placeholder

	mu      sync.Mutex
	running bool
placeholder

var qpsWSCache = &opsWSQPSCache{
	refreshInterval:    qpsWSRefreshInterval,
	requestCountWindow: qpsWSRequestCountWindow,
placeholder

func (c *opsWSQPSCache) start(opsService *service.OpsService) {
	if c == nil || opsService == nil {
		return
placeholder

	for {
		c.mu.Lock()
		if c.running {
			c.mu.Unlock()
			return
	placeholder

		// If a previous refresh loop is currently stopping, wait for it to fully exit.
		done := c.done
		if done != nil {
			c.mu.Unlock()
			<-done

			c.mu.Lock()
			if c.done == done && !c.running {
				c.done = nil
		placeholder
			c.mu.Unlock()
			continue
	placeholder

		c.opsService = opsService
		ctx, cancel := context.WithCancel(context.Background())
		c.cancel = cancel
		c.done = make(chan struct{placeholder)
		done = c.done
		c.running = true
		c.mu.Unlock()

		go func() {
			defer close(done)
			c.refreshLoop(ctx)
	placeholder()
		return
placeholder
placeholder

// Stop stops the background refresh loop.
// It is safe to call multiple times.
func (c *opsWSQPSCache) Stop() {
	if c == nil {
		return
placeholder

	c.mu.Lock()
	if !c.running {
		done := c.done
		c.mu.Unlock()
		if done != nil {
			<-done
	placeholder
		return
placeholder
	cancel := c.cancel
	c.cancel = nil
	c.running = false
	c.opsService = nil
	done := c.done
	c.mu.Unlock()

	if cancel != nil {
		cancel()
placeholder
	if done != nil {
		<-done
placeholder

	c.mu.Lock()
	if c.done == done && !c.running {
		c.done = nil
placeholder
	c.mu.Unlock()
placeholder

func (c *opsWSQPSCache) refreshLoop(ctx context.Context) {
	ticker := time.NewTicker(c.refreshInterval)
	defer ticker.Stop()

	c.refresh(ctx)
	for {
		select {
		case <-ticker.C:
			c.refresh(ctx)
		case <-ctx.Done():
			return
	placeholder
placeholder
placeholder

func (c *opsWSQPSCache) refresh(parentCtx context.Context) {
	if c == nil {
		return
placeholder

	c.mu.Lock()
	opsService := c.opsService
	c.mu.Unlock()
	if opsService == nil {
		return
placeholder

	if parentCtx == nil {
		parentCtx = context.Background()
placeholder
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	now := time.Now().UTC()
	stats, err := opsService.GetWindowStats(ctx, now.Add(-c.requestCountWindow), now)
	if err != nil || stats == nil {
		if err != nil {
			logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] refresh: get window stats failed: %v", err)
	placeholder
		return
placeholder

	requestCount := stats.SuccessCount + stats.ErrorCountTotal
	qps := 0.0
	tps := 0.0
	if c.requestCountWindow > 0 {
		seconds := c.requestCountWindow.Seconds()
		qps = roundTo1DP(float64(requestCount) / seconds)
		tps = roundTo1DP(float64(stats.TokenConsumed) / seconds)
placeholder

	payload := gin.H{
		"type":      "qps_update",
		"timestamp": now.Format(time.RFC3339),
		"data": gin.H{
			"qps":           qps,
			"tps":           tps,
			"request_count": requestCount,
	placeholder,
placeholder

	msg, err := json.Marshal(payload)
	if err != nil {
		logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] refresh: marshal payload failed: %v", err)
		return
placeholder

	c.payload.Store(msg)
	c.lastUpdatedUnixNano.Store(now.UnixNano())
placeholder

func roundTo1DP(v float64) float64 {
	return math.Round(v*10) / 10
placeholder

func (c *opsWSQPSCache) getPayload() []byte {
	if c == nil {
		return nil
placeholder
	if cached, ok := c.payload.Load().([]byte); ok && cached != nil {
		return cached
placeholder
	return nil
placeholder

func closeWS(conn *websocket.Conn, code int, reason string) {
	if conn == nil {
		return
placeholder
	msg := websocket.FormatCloseMessage(code, reason)
	_ = conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(qpsWSWriteTimeout))
	_ = conn.Close()
placeholder

// QPSWSHandler handles realtime QPS push via WebSocket.
// GET /api/v1/admin/ops/ws/qps
func (h *OpsHandler) QPSWSHandler(c *gin.Context) {
	clientIP := requestClientIP(c.Request)

	if h == nil || h.opsService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ops service not initialized"placeholder)
		return
placeholder

	// If realtime monitoring is disabled, prefer a successful WS upgrade followed by a clean close
	// with a deterministic close code. This prevents clients from spinning on 404/1006 reconnect loops.
	if !h.opsService.IsRealtimeMonitoringEnabled(c.Request.Context()) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ops realtime monitoring is disabled"placeholder)
			return
	placeholder
		closeWS(conn, opsWSCloseRealtimeDisabled, "realtime_disabled")
		return
placeholder

	cancelQPSWSIdleStop()
	// Lazily start the background refresh loop so unit tests that never hit the
	// websocket route don't spawn goroutines that depend on DB/Redis stubs.
	qpsWSCache.start(h.opsService)

	// Reserve a global slot before upgrading the connection to keep the limit strict.
	if !tryAcquireOpsWSTotalSlot(opsWSLimits.MaxConns) {
		logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] connection limit reached: %d/%d", wsConnCount.Load(), opsWSLimits.MaxConns)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "too many connections"placeholder)
		return
placeholder
	defer func() {
		if wsConnCount.Add(-1) == 0 {
			scheduleQPSWSIdleStop()
	placeholder
placeholder()

	if opsWSLimits.MaxConnsPerIP > 0 && clientIP != "" {
		if !tryAcquireOpsWSIPSlot(clientIP, opsWSLimits.MaxConnsPerIP) {
			logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] per-ip connection limit reached: ip=%s limit=%d", clientIP, opsWSLimits.MaxConnsPerIP)
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "too many connections"placeholder)
			return
	placeholder
		defer releaseOpsWSIPSlot(clientIP)
placeholder

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] upgrade failed: %v", err)
		return
placeholder

	defer func() {
		_ = conn.Close()
placeholder()

	handleQPSWebSocket(c.Request.Context(), conn)
placeholder

func tryAcquireOpsWSTotalSlot(limit int32) bool {
	if limit <= 0 {
		return true
placeholder
	for {
		current := wsConnCount.Load()
		if current >= limit {
			return false
	placeholder
		if wsConnCount.CompareAndSwap(current, current+1) {
			return true
	placeholder
placeholder
placeholder

func tryAcquireOpsWSIPSlot(clientIP string, limit int32) bool {
	if strings.TrimSpace(clientIP) == "" || limit <= 0 {
		return true
placeholder

	v, _ := wsConnCountByIP.LoadOrStore(clientIP, &atomic.Int32{placeholder)
	counter, ok := v.(*atomic.Int32)
	if !ok {
		return false
placeholder

	for {
		current := counter.Load()
		if current >= limit {
			return false
	placeholder
		if counter.CompareAndSwap(current, current+1) {
			return true
	placeholder
placeholder
placeholder

func releaseOpsWSIPSlot(clientIP string) {
	if strings.TrimSpace(clientIP) == "" {
		return
placeholder

	v, ok := wsConnCountByIP.Load(clientIP)
	if !ok {
		return
placeholder
	counter, ok := v.(*atomic.Int32)
	if !ok {
		return
placeholder
	next := counter.Add(-1)
	if next <= 0 {
		// Best-effort cleanup; safe even if a new slot was acquired concurrently.
		wsConnCountByIP.Delete(clientIP)
placeholder
placeholder

func handleQPSWebSocket(parentCtx context.Context, conn *websocket.Conn) {
	if conn == nil {
		return
placeholder

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	var closeOnce sync.Once
	closeConn := func() {
		closeOnce.Do(func() {
			_ = conn.Close()
	placeholder)
placeholder

	closeFrameCh := make(chan []byte, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		conn.SetReadLimit(qpsWSMaxReadBytes)
		if err := conn.SetReadDeadline(time.Now().Add(qpsWSPongWait)); err != nil {
			logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] set read deadline failed: %v", err)
			return
	placeholder
		conn.SetPongHandler(func(string) error {
			return conn.SetReadDeadline(time.Now().Add(qpsWSPongWait))
	placeholder)
		conn.SetCloseHandler(func(code int, text string) error {
			select {
			case closeFrameCh <- websocket.FormatCloseMessage(code, text):
			default:
		placeholder
			cancel()
			return nil
	placeholder)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
					logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] read failed: %v", err)
			placeholder
				return
		placeholder
	placeholder
placeholder()

	// Push QPS data every 2 seconds (values are globally cached and refreshed at most once per qpsWSRefreshInterval).
	pushTicker := time.NewTicker(qpsWSPushInterval)
	defer pushTicker.Stop()

	// Heartbeat ping every 30 seconds.
	pingTicker := time.NewTicker(qpsWSPingInterval)
	defer pingTicker.Stop()

	writeWithTimeout := func(messageType int, data []byte) error {
		if err := conn.SetWriteDeadline(time.Now().Add(qpsWSWriteTimeout)); err != nil {
			return err
	placeholder
		return conn.WriteMessage(messageType, data)
placeholder

	sendClose := func(closeFrame []byte) {
		if closeFrame == nil {
			closeFrame = websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	placeholder
		_ = writeWithTimeout(websocket.CloseMessage, closeFrame)
placeholder

	for {
		select {
		case <-pushTicker.C:
			msg := qpsWSCache.getPayload()
			if msg == nil {
				continue
		placeholder
			if err := writeWithTimeout(websocket.TextMessage, msg); err != nil {
				logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] write failed: %v", err)
				cancel()
				closeConn()
				wg.Wait()
				return
		placeholder

		case <-pingTicker.C:
			if err := writeWithTimeout(websocket.PingMessage, nil); err != nil {
				logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] ping failed: %v", err)
				cancel()
				closeConn()
				wg.Wait()
				return
		placeholder

		case closeFrame := <-closeFrameCh:
			sendClose(closeFrame)
			closeConn()
			wg.Wait()
			return

		case <-ctx.Done():
			var closeFrame []byte
			select {
			case closeFrame = <-closeFrameCh:
			default:
		placeholder
			sendClose(closeFrame)

			closeConn()
			wg.Wait()
			return
	placeholder
placeholder
placeholder

func isAllowedOpsWSOrigin(r *http.Request) bool {
	if r == nil {
		return false
placeholder
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		switch strings.ToLower(strings.TrimSpace(opsWSProxyConfig.OriginPolicy)) {
		case OriginPolicyStrict:
			return false
		case OriginPolicyPermissive, "":
			return true
		default:
			return true
	placeholder
placeholder
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Hostname() == "" {
		return false
placeholder
	originHost := strings.ToLower(parsed.Hostname())

	trustProxyHeaders := shouldTrustOpsWSProxyHeaders(r)
	reqHost := hostWithoutPort(r.Host)
	if trustProxyHeaders {
		xfHost := strings.TrimSpace(r.Header.Get("X-Forwarded-Host"))
		if xfHost != "" {
			xfHost = strings.TrimSpace(strings.Split(xfHost, ",")[0])
			if xfHost != "" {
				reqHost = hostWithoutPort(xfHost)
		placeholder
	placeholder
placeholder
	reqHost = strings.ToLower(reqHost)
	if reqHost == "" {
		return false
placeholder
	return originHost == reqHost
placeholder

func shouldTrustOpsWSProxyHeaders(r *http.Request) bool {
	if r == nil {
		return false
placeholder
	if !opsWSProxyConfig.TrustProxy {
		return false
placeholder
	peerIP, ok := requestPeerIP(r)
	if !ok {
		return false
placeholder
	return isAddrInTrustedProxies(peerIP, opsWSProxyConfig.TrustedProxies)
placeholder

func requestPeerIP(r *http.Request) (netip.Addr, bool) {
	if r == nil {
		return netip.Addr{placeholder, false
placeholder
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		host = strings.TrimSpace(r.RemoteAddr)
placeholder
	host = strings.TrimPrefix(host, "[")
	host = strings.TrimSuffix(host, "]")
	if host == "" {
		return netip.Addr{placeholder, false
placeholder
	addr, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{placeholder, false
placeholder
	return addr.Unmap(), true
placeholder

func requestClientIP(r *http.Request) string {
	if r == nil {
		return ""
placeholder

	trustProxyHeaders := shouldTrustOpsWSProxyHeaders(r)
	if trustProxyHeaders {
		xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
		if xff != "" {
			// Use the left-most entry (original client). If multiple proxies add values, they are comma-separated.
			xff = strings.TrimSpace(strings.Split(xff, ",")[0])
			xff = strings.TrimPrefix(xff, "[")
			xff = strings.TrimSuffix(xff, "]")
			if addr, err := netip.ParseAddr(xff); err == nil && addr.IsValid() {
				return addr.Unmap().String()
		placeholder
	placeholder
placeholder

	if peer, ok := requestPeerIP(r); ok && peer.IsValid() {
		return peer.String()
placeholder
	return ""
placeholder

func isAddrInTrustedProxies(addr netip.Addr, trusted []netip.Prefix) bool {
	if !addr.IsValid() {
		return false
placeholder
	for _, p := range trusted {
		if p.Contains(addr) {
			return true
	placeholder
placeholder
	return false
placeholder

func loadOpsWSProxyConfigFromEnv() OpsWSProxyConfig {
	cfg := OpsWSProxyConfig{
		TrustProxy:     true,
		TrustedProxies: defaultTrustedProxies(),
		OriginPolicy:   OriginPolicyPermissive,
placeholder

	if v := strings.TrimSpace(os.Getenv(envOpsWSTrustProxy)); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			cfg.TrustProxy = parsed
	placeholder else {
			logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] invalid %s=%q (expected bool); using default=%v", envOpsWSTrustProxy, v, cfg.TrustProxy)
	placeholder
placeholder

	if raw := strings.TrimSpace(os.Getenv(envOpsWSTrustedProxies)); raw != "" {
		prefixes, invalid := parseTrustedProxyList(raw)
		if len(invalid) > 0 {
			logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] invalid %s entries ignored: %s", envOpsWSTrustedProxies, strings.Join(invalid, ", "))
	placeholder
		cfg.TrustedProxies = prefixes
placeholder

	if v := strings.TrimSpace(os.Getenv(envOpsWSOriginPolicy)); v != "" {
		normalized := strings.ToLower(v)
		switch normalized {
		case OriginPolicyStrict, OriginPolicyPermissive:
			cfg.OriginPolicy = normalized
		default:
			logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] invalid %s=%q (expected %q or %q); using default=%q", envOpsWSOriginPolicy, v, OriginPolicyStrict, OriginPolicyPermissive, cfg.OriginPolicy)
	placeholder
placeholder

	return cfg
placeholder

func loadOpsWSRuntimeLimitsFromEnv() opsWSRuntimeLimits {
	cfg := opsWSRuntimeLimits{
		MaxConns:      defaultMaxWSConns,
		MaxConnsPerIP: defaultMaxWSConnsPerIP,
placeholder

	if v := strings.TrimSpace(os.Getenv(envOpsWSMaxConns)); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			cfg.MaxConns = int32(parsed)
	placeholder else {
			logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] invalid %s=%q (expected int>0); using default=%d", envOpsWSMaxConns, v, cfg.MaxConns)
	placeholder
placeholder
	if v := strings.TrimSpace(os.Getenv(envOpsWSMaxConnsPerIP)); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			cfg.MaxConnsPerIP = int32(parsed)
	placeholder else {
			logger.LegacyPrintf("handler.admin.ops_ws", "[OpsWS] invalid %s=%q (expected int>=0); using default=%d", envOpsWSMaxConnsPerIP, v, cfg.MaxConnsPerIP)
	placeholder
placeholder
	return cfg
placeholder

func defaultTrustedProxies() []netip.Prefix {
	prefixes, _ := parseTrustedProxyList("127.0.0.0/8,::1/128")
	return prefixes
placeholder

func parseTrustedProxyList(raw string) (prefixes []netip.Prefix, invalid []string) {
	for _, token := range strings.Split(raw, ",") {
		item := strings.TrimSpace(token)
		if item == "" {
			continue
	placeholder

		var (
			p   netip.Prefix
			err error
		)
		if strings.Contains(item, "/") {
			p, err = netip.ParsePrefix(item)
	placeholder else {
			var addr netip.Addr
			addr, err = netip.ParseAddr(item)
			if err == nil {
				addr = addr.Unmap()
				bits := 128
				if addr.Is4() {
					bits = 32
			placeholder
				p = netip.PrefixFrom(addr, bits)
		placeholder
	placeholder

		if err != nil || !p.IsValid() {
			invalid = append(invalid, item)
			continue
	placeholder

		prefixes = append(prefixes, p.Masked())
placeholder
	return prefixes, invalid
placeholder

func hostWithoutPort(hostport string) string {
	hostport = strings.TrimSpace(hostport)
	if hostport == "" {
		return ""
placeholder
	if host, _, err := net.SplitHostPort(hostport); err == nil {
		return host
placeholder
	if strings.HasPrefix(hostport, "[") && strings.HasSuffix(hostport, "]") {
		return strings.Trim(hostport, "[]")
placeholder
	parts := strings.Split(hostport, ":")
	return parts[0]
placeholder
