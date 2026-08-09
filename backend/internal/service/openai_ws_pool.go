package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"golang.org/x/sync/errgroup"
)

const (
	openAIWSConnMaxAge             = 60 * time.Minute
	openAIWSConnHealthCheckIdle    = 90 * time.Second
	openAIWSConnHealthCheckTO      = 2 * time.Second
	openAIWSConnPrewarmExtraDelay  = 2 * time.Second
	openAIWSAcquireCleanupInterval = 3 * time.Second
	openAIWSBackgroundPingInterval = 30 * time.Second
	openAIWSBackgroundSweepTicker  = 30 * time.Second

	openAIWSPrewarmFailureWindow   = 30 * time.Second
	openAIWSPrewarmFailureSuppress = 2
)

var (
	errOpenAIWSConnClosed               = errors.New("openai ws connection closed")
	errOpenAIWSConnQueueFull            = errors.New("openai ws connection queue full")
	errOpenAIWSPreferredConnUnavailable = errors.New("openai ws preferred connection unavailable")
)

type openAIWSDialError struct {
	StatusCode      int
	ResponseHeaders http.Header
	ResponseBody    []byte
	Err             error
placeholder

func (e *openAIWSDialError) Error() string {
	if e == nil {
		return ""
placeholder
	if e.StatusCode > 0 {
		return fmt.Sprintf("openai ws dial failed: status=%d err=%v", e.StatusCode, e.Err)
placeholder
	return fmt.Sprintf("openai ws dial failed: %v", e.Err)
placeholder

func (e *openAIWSDialError) Unwrap() error {
	if e == nil {
		return nil
placeholder
	return e.Err
placeholder

type openAIWSAcquireRequest struct {
	Account *Account
	WSURL   string
	Headers http.Header
	// HeadersFactory is evaluated inside dialConn. It exists so credentials
	// whose authorization is per-dial (Agent Identity) are never cached in
	// lastAcquire or delayed prewarm state.
	HeadersFactory  func(context.Context, http.Header) (http.Header, error)
	ProxyURL        string
	PreferredConnID string
	// ForceNewConn: 强制本次获取新连接（避免复用导致连接内续链状态互相污染）。
	ForceNewConn bool
	// ForcePreferredConn: 强制本次只使用 PreferredConnID，禁止漂移到其它连接。
	ForcePreferredConn bool
placeholder

type openAIWSHandshakeCompatibilityKey struct {
	betaFeatures string
placeholder

type openAIWSConnLease struct {
	pool      *openAIWSConnPool
	accountID int64
	conn      *openAIWSConn
	queueWait time.Duration
	connPick  time.Duration
	reused    bool
	released  atomic.Bool
placeholder

func (l *openAIWSConnLease) activeConn() (*openAIWSConn, error) {
	if l == nil || l.conn == nil {
		return nil, errOpenAIWSConnClosed
placeholder
	if l.released.Load() {
		return nil, errOpenAIWSConnClosed
placeholder
	return l.conn, nil
placeholder

func (l *openAIWSConnLease) ConnID() string {
	if l == nil || l.conn == nil {
		return ""
placeholder
	return l.conn.id
placeholder

func (l *openAIWSConnLease) QueueWaitDuration() time.Duration {
	if l == nil {
		return 0
placeholder
	return l.queueWait
placeholder

func (l *openAIWSConnLease) ConnPickDuration() time.Duration {
	if l == nil {
		return 0
placeholder
	return l.connPick
placeholder

func (l *openAIWSConnLease) Reused() bool {
	if l == nil {
		return false
placeholder
	return l.reused
placeholder

func (l *openAIWSConnLease) HandshakeHeader(name string) string {
	if l == nil || l.conn == nil {
		return ""
placeholder
	return l.conn.handshakeHeader(name)
placeholder

func (l *openAIWSConnLease) HandshakeHeaders() http.Header {
	if l == nil || l.conn == nil {
		return nil
placeholder
	return cloneHeader(l.conn.handshakeHeaders)
placeholder

func (l *openAIWSConnLease) IsPrewarmed() bool {
	if l == nil || l.conn == nil {
		return false
placeholder
	return l.conn.isPrewarmed()
placeholder

func (l *openAIWSConnLease) MarkPrewarmed() {
	if l == nil || l.conn == nil {
		return
placeholder
	l.conn.markPrewarmed()
placeholder

func (l *openAIWSConnLease) WriteJSON(value any, timeout time.Duration) error {
	conn, err := l.activeConn()
	if err != nil {
		return err
placeholder
	return conn.writeJSONWithTimeout(context.Background(), value, timeout)
placeholder

func (l *openAIWSConnLease) WriteJSONWithContextTimeout(ctx context.Context, value any, timeout time.Duration) error {
	conn, err := l.activeConn()
	if err != nil {
		return err
placeholder
	return conn.writeJSONWithTimeout(ctx, value, timeout)
placeholder

func (l *openAIWSConnLease) WriteJSONContext(ctx context.Context, value any) error {
	conn, err := l.activeConn()
	if err != nil {
		return err
placeholder
	return conn.writeJSON(value, ctx)
placeholder

func (l *openAIWSConnLease) ReadMessage(timeout time.Duration) ([]byte, error) {
	conn, err := l.activeConn()
	if err != nil {
		return nil, err
placeholder
	return conn.readMessageWithTimeout(timeout)
placeholder

func (l *openAIWSConnLease) ReadMessageContext(ctx context.Context) ([]byte, error) {
	conn, err := l.activeConn()
	if err != nil {
		return nil, err
placeholder
	return conn.readMessage(ctx)
placeholder

func (l *openAIWSConnLease) ReadMessageWithContextTimeout(ctx context.Context, timeout time.Duration) ([]byte, error) {
	conn, err := l.activeConn()
	if err != nil {
		return nil, err
placeholder
	return conn.readMessageWithContextTimeout(ctx, timeout)
placeholder

func (l *openAIWSConnLease) PingWithTimeout(timeout time.Duration) error {
	conn, err := l.activeConn()
	if err != nil {
		return err
placeholder
	return conn.pingWithTimeout(timeout)
placeholder

func (l *openAIWSConnLease) SupportsIdlePingWithoutReader() bool {
	conn, err := l.activeConn()
	if err != nil {
		return false
placeholder
	return conn.supportsIdlePingWithoutReader()
placeholder

func (l *openAIWSConnLease) MarkBroken() {
	if l == nil || l.pool == nil || l.conn == nil || l.released.Load() {
		return
placeholder
	l.pool.evictConn(l.accountID, l.conn.id)
placeholder

func (l *openAIWSConnLease) Release() {
	if l == nil || l.conn == nil {
		return
placeholder
	if !l.released.CompareAndSwap(false, true) {
		return
placeholder
	l.conn.release()
	if l.pool != nil {
		l.pool.notifyAccountPoolChanged(l.accountID)
placeholder
placeholder

type openAIWSConn struct {
	id string
	ws openAIWSClientConn

	handshakeHeaders       http.Header
	handshakeCompatibility openAIWSHandshakeCompatibilityKey
	routingAffinity        string

	leaseCh   chan struct{placeholder
	closedCh  chan struct{placeholder
	closeOnce sync.Once

	readMu  sync.Mutex
	writeMu sync.Mutex

	waiters       atomic.Int32
	createdAtNano atomic.Int64
	lastUsedNano  atomic.Int64
	prewarmed     atomic.Bool
placeholder

func newOpenAIWSConn(id string, _ int64, ws openAIWSClientConn, handshakeHeaders http.Header) *openAIWSConn {
	now := time.Now()
	conn := &openAIWSConn{
		id:               id,
		ws:               ws,
		handshakeHeaders: cloneHeader(handshakeHeaders),
		leaseCh:          make(chan struct{placeholder, 1),
		closedCh:         make(chan struct{placeholder),
placeholder
	conn.leaseCh <- struct{placeholder{placeholder
	conn.createdAtNano.Store(now.UnixNano())
	conn.lastUsedNano.Store(now.UnixNano())
	return conn
placeholder

func (c *openAIWSConn) tryAcquire() bool {
	if c == nil {
		return false
placeholder
	select {
	case <-c.closedCh:
		return false
	default:
placeholder
	select {
	case <-c.leaseCh:
		select {
		case <-c.closedCh:
			c.release()
			return false
		default:
	placeholder
		return true
	default:
		return false
placeholder
placeholder

func (c *openAIWSConn) acquire(ctx context.Context) error {
	if c == nil {
		return errOpenAIWSConnClosed
placeholder
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.closedCh:
			return errOpenAIWSConnClosed
		case <-c.leaseCh:
			// A cancellation and a lease delivery can become ready together. Once
			// the semaphore token has been consumed, check the context again and
			// return it before reporting cancellation so a canceled waiter cannot
			// strand a pooled connection.
			if err := ctx.Err(); err != nil {
				c.release()
				return err
		placeholder
			select {
			case <-c.closedCh:
				c.release()
				return errOpenAIWSConnClosed
			default:
		placeholder
			return nil
	placeholder
placeholder
placeholder

func (c *openAIWSConn) release() {
	if c == nil {
		return
placeholder
	select {
	case c.leaseCh <- struct{placeholder{placeholder:
	default:
placeholder
	c.touch()
placeholder

func (c *openAIWSConn) close() {
	if c == nil {
		return
placeholder
	c.closeOnce.Do(func() {
		close(c.closedCh)
		if c.ws != nil {
			_ = c.ws.Close()
	placeholder
		select {
		case c.leaseCh <- struct{placeholder{placeholder:
		default:
	placeholder
placeholder)
placeholder

func (c *openAIWSConn) writeJSONWithTimeout(parent context.Context, value any, timeout time.Duration) error {
	if c == nil {
		return errOpenAIWSConnClosed
placeholder
	select {
	case <-c.closedCh:
		return errOpenAIWSConnClosed
	default:
placeholder

	writeCtx := parent
	if writeCtx == nil {
		writeCtx = context.Background()
placeholder
	if timeout <= 0 {
		return c.writeJSON(value, writeCtx)
placeholder
	var cancel context.CancelFunc
	writeCtx, cancel = context.WithTimeout(writeCtx, timeout)
	defer cancel()
	return c.writeJSON(value, writeCtx)
placeholder

func (c *openAIWSConn) writeJSON(value any, writeCtx context.Context) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.ws == nil {
		return errOpenAIWSConnClosed
placeholder
	if writeCtx == nil {
		writeCtx = context.Background()
placeholder
	if err := c.ws.WriteJSON(writeCtx, value); err != nil {
		return err
placeholder
	c.touch()
	return nil
placeholder

func (c *openAIWSConn) readMessageWithTimeout(timeout time.Duration) ([]byte, error) {
	return c.readMessageWithContextTimeout(context.Background(), timeout)
placeholder

func (c *openAIWSConn) readMessageWithContextTimeout(parent context.Context, timeout time.Duration) ([]byte, error) {
	if c == nil {
		return nil, errOpenAIWSConnClosed
placeholder
	select {
	case <-c.closedCh:
		return nil, errOpenAIWSConnClosed
	default:
placeholder

	if parent == nil {
		parent = context.Background()
placeholder
	if timeout <= 0 {
		return c.readMessage(parent)
placeholder
	readCtx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	return c.readMessage(readCtx)
placeholder

func (c *openAIWSConn) readMessage(readCtx context.Context) ([]byte, error) {
	c.readMu.Lock()
	defer c.readMu.Unlock()
	if c.ws == nil {
		return nil, errOpenAIWSConnClosed
placeholder
	if readCtx == nil {
		readCtx = context.Background()
placeholder
	payload, err := c.ws.ReadMessage(readCtx)
	if err != nil {
		return nil, err
placeholder
	c.touch()
	return payload, nil
placeholder

func (c *openAIWSConn) pingWithTimeout(timeout time.Duration) error {
	if c == nil {
		return errOpenAIWSConnClosed
placeholder
	select {
	case <-c.closedCh:
		return errOpenAIWSConnClosed
	default:
placeholder

	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.ws == nil {
		return errOpenAIWSConnClosed
placeholder
	if timeout <= 0 {
		timeout = openAIWSConnHealthCheckTO
placeholder
	pingCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := c.ws.Ping(pingCtx); err != nil {
		return err
placeholder
	return nil
placeholder

func (c *openAIWSConn) supportsIdlePingWithoutReader() bool {
	if c == nil || c.ws == nil {
		return false
placeholder
	capable, ok := c.ws.(openAIWSIdlePingCapable)
	// Test and alternate implementations keep the historical probe behavior
	// unless they explicitly declare it unsafe.
	return !ok || capable.SupportsIdlePingWithoutReader()
placeholder

func (c *openAIWSConn) touch() {
	if c == nil {
		return
placeholder
	c.lastUsedNano.Store(time.Now().UnixNano())
placeholder

func (c *openAIWSConn) createdAt() time.Time {
	if c == nil {
		return time.Time{placeholder
placeholder
	nano := c.createdAtNano.Load()
	if nano <= 0 {
		return time.Time{placeholder
placeholder
	return time.Unix(0, nano)
placeholder

func (c *openAIWSConn) lastUsedAt() time.Time {
	if c == nil {
		return time.Time{placeholder
placeholder
	nano := c.lastUsedNano.Load()
	if nano <= 0 {
		return time.Time{placeholder
placeholder
	return time.Unix(0, nano)
placeholder

func (c *openAIWSConn) idleDuration(now time.Time) time.Duration {
	if c == nil {
		return 0
placeholder
	last := c.lastUsedAt()
	if last.IsZero() {
		return 0
placeholder
	return now.Sub(last)
placeholder

func (c *openAIWSConn) age(now time.Time) time.Duration {
	if c == nil {
		return 0
placeholder
	created := c.createdAt()
	if created.IsZero() {
		return 0
placeholder
	return now.Sub(created)
placeholder

func (c *openAIWSConn) isLeased() bool {
	if c == nil {
		return false
placeholder
	return len(c.leaseCh) == 0
placeholder

func (c *openAIWSConn) handshakeHeader(name string) string {
	if c == nil || c.handshakeHeaders == nil {
		return ""
placeholder
	return strings.TrimSpace(c.handshakeHeaders.Get(strings.TrimSpace(name)))
placeholder

func (c *openAIWSConn) matchesHandshakeCompatibility(compatibility openAIWSHandshakeCompatibilityKey) bool {
	return c != nil && c.handshakeCompatibility == compatibility
placeholder

func (c *openAIWSConn) matchesRoutingAffinity(routingAffinity string) bool {
	return c != nil && c.routingAffinity == routingAffinity
placeholder

func (c *openAIWSConn) isPrewarmed() bool {
	if c == nil {
		return false
placeholder
	return c.prewarmed.Load()
placeholder

func (c *openAIWSConn) markPrewarmed() {
	if c == nil {
		return
placeholder
	c.prewarmed.Store(true)
placeholder

type openAIWSAccountPool struct {
	mu            sync.Mutex
	conns         map[string]*openAIWSConn
	pinnedConns   map[string]int
	changedCh     chan struct{placeholder
	creating      int
	generation    uint64
	lastCleanupAt time.Time
	lastAcquire   *openAIWSAcquireRequest
	prewarmActive bool
	prewarmUntil  time.Time
	prewarmFails  int
	prewarmFailAt time.Time
placeholder

func (ap *openAIWSAccountPool) changeChannelLocked() chan struct{placeholder {
	if ap.changedCh == nil {
		ap.changedCh = make(chan struct{placeholder)
placeholder
	return ap.changedCh
placeholder

func (ap *openAIWSAccountPool) signalChangedLocked() {
	if ap == nil {
		return
placeholder
	if ap.changedCh != nil {
		close(ap.changedCh)
placeholder
	ap.changedCh = make(chan struct{placeholder)
placeholder

type OpenAIWSPoolMetricsSnapshot struct {
	AcquireTotal            int64
	AcquireReuseTotal       int64
	AcquireCreateTotal      int64
	AcquireQueueWaitTotal   int64
	AcquireQueueWaitMsTotal int64
	ConnPickTotal           int64
	ConnPickMsTotal         int64
	ScaleUpTotal            int64
	ScaleDownTotal          int64
placeholder

type openAIWSPoolMetrics struct {
	acquireTotal          atomic.Int64
	acquireReuseTotal     atomic.Int64
	acquireCreateTotal    atomic.Int64
	acquireQueueWaitTotal atomic.Int64
	acquireQueueWaitMs    atomic.Int64
	connPickTotal         atomic.Int64
	connPickMs            atomic.Int64
	scaleUpTotal          atomic.Int64
	scaleDownTotal        atomic.Int64
placeholder

type openAIWSConnPool struct {
	cfg *config.Config
	// 通过接口解耦底层 WS 客户端实现，默认使用 coder/websocket。
	clientDialer openAIWSClientDialer

	accounts sync.Map // key: int64(accountID), value: *openAIWSAccountPool
	seq      atomic.Uint64

	metrics openAIWSPoolMetrics

	workerStopCh chan struct{placeholder
	workerWg     sync.WaitGroup
	closeOnce    sync.Once
placeholder

func newOpenAIWSConnPool(cfg *config.Config) *openAIWSConnPool {
	pool := &openAIWSConnPool{
		cfg:          cfg,
		clientDialer: newDefaultOpenAIWSClientDialer(),
		workerStopCh: make(chan struct{placeholder),
placeholder
	pool.startBackgroundWorkers()
	return pool
placeholder

func (p *openAIWSConnPool) SnapshotMetrics() OpenAIWSPoolMetricsSnapshot {
	if p == nil {
		return OpenAIWSPoolMetricsSnapshot{placeholder
placeholder
	return OpenAIWSPoolMetricsSnapshot{
		AcquireTotal:            p.metrics.acquireTotal.Load(),
		AcquireReuseTotal:       p.metrics.acquireReuseTotal.Load(),
		AcquireCreateTotal:      p.metrics.acquireCreateTotal.Load(),
		AcquireQueueWaitTotal:   p.metrics.acquireQueueWaitTotal.Load(),
		AcquireQueueWaitMsTotal: p.metrics.acquireQueueWaitMs.Load(),
		ConnPickTotal:           p.metrics.connPickTotal.Load(),
		ConnPickMsTotal:         p.metrics.connPickMs.Load(),
		ScaleUpTotal:            p.metrics.scaleUpTotal.Load(),
		ScaleDownTotal:          p.metrics.scaleDownTotal.Load(),
placeholder
placeholder

func (p *openAIWSConnPool) SnapshotTransportMetrics() OpenAIWSTransportMetricsSnapshot {
	if p == nil {
		return OpenAIWSTransportMetricsSnapshot{placeholder
placeholder
	if dialer, ok := p.clientDialer.(openAIWSTransportMetricsDialer); ok {
		return dialer.SnapshotTransportMetrics()
placeholder
	return OpenAIWSTransportMetricsSnapshot{placeholder
placeholder

func (p *openAIWSConnPool) setClientDialerForTest(dialer openAIWSClientDialer) {
	if p == nil || dialer == nil {
		return
placeholder
	p.clientDialer = dialer
placeholder

// Close 停止后台 worker 并关闭所有空闲连接，应在优雅关闭时调用。
func (p *openAIWSConnPool) Close() {
	if p == nil {
		return
placeholder
	p.closeOnce.Do(func() {
		if p.workerStopCh != nil {
			close(p.workerStopCh)
	placeholder
		p.workerWg.Wait()
		// 遍历所有账户池，关闭全部空闲连接。
		p.accounts.Range(func(key, value any) bool {
			ap, ok := value.(*openAIWSAccountPool)
			if !ok || ap == nil {
				return true
		placeholder
			ap.mu.Lock()
			for _, conn := range ap.conns {
				if conn != nil && !conn.isLeased() {
					conn.close()
			placeholder
		placeholder
			ap.mu.Unlock()
			return true
	placeholder)
placeholder)
placeholder

func (p *openAIWSConnPool) startBackgroundWorkers() {
	if p == nil || p.workerStopCh == nil {
		return
placeholder
	p.workerWg.Add(2)
	go func() {
		defer p.workerWg.Done()
		p.runBackgroundPingWorker()
placeholder()
	go func() {
		defer p.workerWg.Done()
		p.runBackgroundCleanupWorker()
placeholder()
placeholder

type openAIWSIdlePingCandidate struct {
	accountID int64
	conn      *openAIWSConn
placeholder

func (p *openAIWSConnPool) runBackgroundPingWorker() {
	if p == nil {
		return
placeholder
	ticker := time.NewTicker(openAIWSBackgroundPingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			p.runBackgroundPingSweep()
		case <-p.workerStopCh:
			return
	placeholder
placeholder
placeholder

func (p *openAIWSConnPool) runBackgroundPingSweep() {
	if p == nil {
		return
placeholder
	candidates := p.snapshotIdleConnsForPing()
	var g errgroup.Group
	g.SetLimit(10)
	for _, item := range candidates {
		item := item
		if item.conn == nil || item.conn.isLeased() || item.conn.waiters.Load() > 0 || !item.conn.supportsIdlePingWithoutReader() {
			continue
	placeholder
		g.Go(func() error {
			if err := item.conn.pingWithTimeout(openAIWSConnHealthCheckTO); err != nil {
				p.evictConn(item.accountID, item.conn.id)
		placeholder
			return nil
	placeholder)
placeholder
	_ = g.Wait()
placeholder

func (p *openAIWSConnPool) snapshotIdleConnsForPing() []openAIWSIdlePingCandidate {
	if p == nil {
		return nil
placeholder
	candidates := make([]openAIWSIdlePingCandidate, 0)
	p.accounts.Range(func(key, value any) bool {
		accountID, ok := key.(int64)
		if !ok || accountID <= 0 {
			return true
	placeholder
		ap, ok := value.(*openAIWSAccountPool)
		if !ok || ap == nil {
			return true
	placeholder
		ap.mu.Lock()
		for _, conn := range ap.conns {
			if conn == nil || conn.isLeased() || conn.waiters.Load() > 0 {
				continue
		placeholder
			candidates = append(candidates, openAIWSIdlePingCandidate{
				accountID: accountID,
				conn:      conn,
		placeholder)
	placeholder
		ap.mu.Unlock()
		return true
placeholder)
	return candidates
placeholder

func (p *openAIWSConnPool) runBackgroundCleanupWorker() {
	if p == nil {
		return
placeholder
	ticker := time.NewTicker(openAIWSBackgroundSweepTicker)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			p.runBackgroundCleanupSweep(time.Now())
		case <-p.workerStopCh:
			return
	placeholder
placeholder
placeholder

func (p *openAIWSConnPool) runBackgroundCleanupSweep(now time.Time) {
	if p == nil {
		return
placeholder
	type cleanupResult struct {
		evicted []*openAIWSConn
placeholder
	results := make([]cleanupResult, 0)
	p.accounts.Range(func(_ any, value any) bool {
		ap, ok := value.(*openAIWSAccountPool)
		if !ok || ap == nil {
			return true
	placeholder
		maxConns := p.maxConnsHardCap()
		ap.mu.Lock()
		if ap.lastAcquire != nil && ap.lastAcquire.Account != nil {
			maxConns = p.effectiveMaxConnsByAccount(ap.lastAcquire.Account)
	placeholder
		evicted := p.cleanupAccountLocked(ap, now, maxConns)
		ap.lastCleanupAt = now
		ap.mu.Unlock()
		if len(evicted) > 0 {
			results = append(results, cleanupResult{evicted: evictedplaceholder)
	placeholder
		return true
placeholder)
	for _, result := range results {
		closeOpenAIWSConns(result.evicted)
placeholder
placeholder

func (p *openAIWSConnPool) Acquire(ctx context.Context, req openAIWSAcquireRequest) (*openAIWSConnLease, error) {
	if p != nil {
		p.metrics.acquireTotal.Add(1)
placeholder
	return p.acquire(ctx, cloneOpenAIWSAcquireRequest(req), 0)
placeholder

func (p *openAIWSConnPool) acquire(ctx context.Context, req openAIWSAcquireRequest, retry int) (*openAIWSConnLease, error) {
	if p == nil || req.Account == nil || req.Account.ID <= 0 {
		return nil, errors.New("invalid ws acquire request")
placeholder
	if stringsTrim(req.WSURL) == "" {
		return nil, errors.New("ws url is empty")
placeholder

retryAcquire:
	accountID := req.Account.ID
	compatibility := normalizeOpenAIWSHandshakeCompatibility(req.Headers)
	routingAffinity := normalizeOpenAIWSRoutingAffinity(req.Headers)
	effectiveMaxConns := p.effectiveMaxConnsByAccount(req.Account)
	if effectiveMaxConns <= 0 {
		return nil, errOpenAIWSConnQueueFull
placeholder
	var evicted []*openAIWSConn
	ap := p.getOrCreateAccountPool(accountID)
	ap.mu.Lock()
	acquireGeneration := ap.generation
	now := time.Now()
	if ap.lastCleanupAt.IsZero() || now.Sub(ap.lastCleanupAt) >= openAIWSAcquireCleanupInterval {
		evicted = p.cleanupAccountLocked(ap, now, effectiveMaxConns)
		ap.lastCleanupAt = now
placeholder
	pickStartedAt := time.Now()
	allowReuse := !req.ForceNewConn
	preferredConnID := stringsTrim(req.PreferredConnID)
	forcePreferredConn := allowReuse && req.ForcePreferredConn

	if allowReuse {
		if forcePreferredConn {
			if preferredConnID == "" {
				p.recordConnPickDuration(time.Since(pickStartedAt))
				ap.mu.Unlock()
				closeOpenAIWSConns(evicted)
				return nil, errOpenAIWSPreferredConnUnavailable
		placeholder
			preferredConn, ok := ap.conns[preferredConnID]
			if !ok || !preferredConn.matchesHandshakeCompatibility(compatibility) {
				p.recordConnPickDuration(time.Since(pickStartedAt))
				ap.mu.Unlock()
				closeOpenAIWSConns(evicted)
				return nil, errOpenAIWSPreferredConnUnavailable
		placeholder
			if preferredConn.tryAcquire() {
				connPick := time.Since(pickStartedAt)
				p.recordConnPickDuration(connPick)
				ap.mu.Unlock()
				closeOpenAIWSConns(evicted)
				if p.shouldHealthCheckConn(preferredConn) {
					if err := preferredConn.pingWithTimeout(openAIWSConnHealthCheckTO); err != nil {
						preferredConn.close()
						p.evictConn(accountID, preferredConn.id)
						if retry < 1 {
							return p.acquire(ctx, req, retry+1)
					placeholder
						return nil, err
				placeholder
			placeholder
				lease := &openAIWSConnLease{
					pool:      p,
					accountID: accountID,
					conn:      preferredConn,
					connPick:  connPick,
					reused:    true,
			placeholder
				p.metrics.acquireReuseTotal.Add(1)
				p.recordLastSuccessfulAcquire(accountID, acquireGeneration, req)
				p.ensureTargetIdleAsync(accountID)
				return lease, nil
		placeholder

			connPick := time.Since(pickStartedAt)
			p.recordConnPickDuration(connPick)
			if int(preferredConn.waiters.Load()) >= p.queueLimitPerConn() {
				ap.mu.Unlock()
				closeOpenAIWSConns(evicted)
				return nil, errOpenAIWSConnQueueFull
		placeholder
			preferredConn.waiters.Add(1)
			ap.mu.Unlock()
			closeOpenAIWSConns(evicted)
			defer preferredConn.waiters.Add(-1)
			waitStart := time.Now()
			p.metrics.acquireQueueWaitTotal.Add(1)

			if err := preferredConn.acquire(ctx); err != nil {
				if errors.Is(err, errOpenAIWSConnClosed) && retry < 1 {
					return p.acquire(ctx, req, retry+1)
			placeholder
				return nil, err
		placeholder
			if p.shouldHealthCheckConn(preferredConn) {
				if err := preferredConn.pingWithTimeout(openAIWSConnHealthCheckTO); err != nil {
					preferredConn.release()
					preferredConn.close()
					p.evictConn(accountID, preferredConn.id)
					if retry < 1 {
						return p.acquire(ctx, req, retry+1)
				placeholder
					return nil, err
			placeholder
		placeholder

			queueWait := time.Since(waitStart)
			p.metrics.acquireQueueWaitMs.Add(queueWait.Milliseconds())
			lease := &openAIWSConnLease{
				pool:      p,
				accountID: accountID,
				conn:      preferredConn,
				queueWait: queueWait,
				connPick:  connPick,
				reused:    true,
		placeholder
			p.metrics.acquireReuseTotal.Add(1)
			p.recordLastSuccessfulAcquire(accountID, acquireGeneration, req)
			p.ensureTargetIdleAsync(accountID)
			return lease, nil
	placeholder

		if preferredConnID != "" {
			if conn, ok := ap.conns[preferredConnID]; ok && conn.matchesHandshakeCompatibility(compatibility) && conn.tryAcquire() {
				connPick := time.Since(pickStartedAt)
				p.recordConnPickDuration(connPick)
				ap.mu.Unlock()
				closeOpenAIWSConns(evicted)
				if p.shouldHealthCheckConn(conn) {
					if err := conn.pingWithTimeout(openAIWSConnHealthCheckTO); err != nil {
						conn.close()
						p.evictConn(accountID, conn.id)
						if retry < 1 {
							return p.acquire(ctx, req, retry+1)
					placeholder
						return nil, err
				placeholder
			placeholder
				lease := &openAIWSConnLease{pool: p, accountID: accountID, conn: conn, connPick: connPick, reused: trueplaceholder
				p.metrics.acquireReuseTotal.Add(1)
				p.recordLastSuccessfulAcquire(accountID, acquireGeneration, req)
				p.ensureTargetIdleAsync(accountID)
				return lease, nil
		placeholder
	placeholder

		// A routing hint is advisory at WebSocket dial time. Prefer a pooled
		// connection whose handshake used the same hint, but do not make that
		// preference a continuation compatibility requirement.
		best := p.pickLeastBusyConnWithRoutingAffinityLocked(ap, compatibility, routingAffinity)
		if best != nil && best.tryAcquire() {
			connPick := time.Since(pickStartedAt)
			p.recordConnPickDuration(connPick)
			ap.mu.Unlock()
			closeOpenAIWSConns(evicted)
			if p.shouldHealthCheckConn(best) {
				if err := best.pingWithTimeout(openAIWSConnHealthCheckTO); err != nil {
					best.close()
					p.evictConn(accountID, best.id)
					if retry < 1 {
						return p.acquire(ctx, req, retry+1)
				placeholder
					return nil, err
			placeholder
		placeholder
			lease := &openAIWSConnLease{pool: p, accountID: accountID, conn: best, connPick: connPick, reused: trueplaceholder
			p.metrics.acquireReuseTotal.Add(1)
			p.recordLastSuccessfulAcquire(accountID, acquireGeneration, req)
			p.ensureTargetIdleAsync(accountID)
			return lease, nil
	placeholder
		for _, conn := range ap.conns {
			if conn == nil || conn == best || !conn.matchesHandshakeCompatibility(compatibility) || !conn.matchesRoutingAffinity(routingAffinity) {
				continue
		placeholder
			if conn.tryAcquire() {
				connPick := time.Since(pickStartedAt)
				p.recordConnPickDuration(connPick)
				ap.mu.Unlock()
				closeOpenAIWSConns(evicted)
				if p.shouldHealthCheckConn(conn) {
					if err := conn.pingWithTimeout(openAIWSConnHealthCheckTO); err != nil {
						conn.close()
						p.evictConn(accountID, conn.id)
						if retry < 1 {
							return p.acquire(ctx, req, retry+1)
					placeholder
						return nil, err
				placeholder
			placeholder
				lease := &openAIWSConnLease{pool: p, accountID: accountID, conn: conn, connPick: connPick, reused: trueplaceholder
				p.metrics.acquireReuseTotal.Add(1)
				p.recordLastSuccessfulAcquire(accountID, acquireGeneration, req)
				p.ensureTargetIdleAsync(accountID)
				return lease, nil
		placeholder
	placeholder
placeholder

	if !req.ForceNewConn && len(ap.conns)+ap.creating >= effectiveMaxConns {
		affine := p.pickLeastBusyConnWithRoutingAffinityLocked(ap, compatibility, routingAffinity)
		if idle := p.pickOldestIdleConnWithoutHandshakeCompatibilityOrRoutingAffinityLocked(ap, compatibility, routingAffinity); idle != nil {
			delete(ap.conns, idle.id)
			evicted = append(evicted, idle)
			p.metrics.scaleDownTotal.Add(1)
	placeholder else if affine == nil {
			compatible := p.pickLeastBusyConnLocked(ap, "", compatibility)
			if compatible != nil {
				// Capacity is full and every compatible connection is busy. The
				// hint remains soft here: queue on a compatible connection below.
				goto acquireAtCapacity
		placeholder
			hasConnection := false
			for _, conn := range ap.conns {
				if conn != nil {
					hasConnection = true
					break
			placeholder
		placeholder
			if !hasConnection && ap.creating == 0 {
				ap.mu.Unlock()
				closeOpenAIWSConns(evicted)
				return nil, errOpenAIWSConnClosed
		placeholder
			changedCh := ap.changeChannelLocked()
			ap.mu.Unlock()
			closeOpenAIWSConns(evicted)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-changedCh:
				goto retryAcquire
		placeholder
	placeholder
placeholder

	if req.ForceNewConn && len(ap.conns)+ap.creating >= effectiveMaxConns {
		if idle := p.pickOldestIdleConnLocked(ap); idle != nil {
			delete(ap.conns, idle.id)
			evicted = append(evicted, idle)
			p.metrics.scaleDownTotal.Add(1)
	placeholder
placeholder

	if len(ap.conns)+ap.creating < effectiveMaxConns {
		connPick := time.Since(pickStartedAt)
		p.recordConnPickDuration(connPick)
		ap.creating++
		ap.mu.Unlock()
		closeOpenAIWSConns(evicted)

		conn, dialErr := p.dialConn(ctx, req)

		ap = p.getOrCreateAccountPool(accountID)
		ap.mu.Lock()
		ap.creating--
		if ap.generation != acquireGeneration {
			ap.signalChangedLocked()
			ap.mu.Unlock()
			if conn != nil {
				conn.close()
		placeholder
			if retry < 1 {
				return p.acquire(ctx, req, retry+1)
		placeholder
			return nil, errOpenAIWSConnClosed
	placeholder
		if dialErr != nil {
			ap.prewarmFails++
			ap.prewarmFailAt = time.Now()
			ap.signalChangedLocked()
			ap.mu.Unlock()
			return nil, dialErr
	placeholder
		// Claim the freshly dialed connection before publishing it. Otherwise a
		// topology waiter awakened below can take the free semaphore first and
		// make the caller that paid for the dial queue behind it.
		if !conn.tryAcquire() {
			ap.signalChangedLocked()
			ap.mu.Unlock()
			conn.close()
			return nil, errOpenAIWSConnClosed
	placeholder
		ap.conns[conn.id] = conn
		ap.prewarmFails = 0
		ap.prewarmFailAt = time.Time{placeholder
		// Wake acquires that observed creating>0 with no compatible connection.
		// Without this signal they can remain asleep until the new lease is
		// released, even though the pool topology already changed.
		ap.signalChangedLocked()
		ap.mu.Unlock()
		p.metrics.acquireCreateTotal.Add(1)
		lease := &openAIWSConnLease{pool: p, accountID: accountID, conn: conn, connPick: connPickplaceholder
		p.recordLastSuccessfulAcquire(accountID, acquireGeneration, req)
		p.ensureTargetIdleAsync(accountID)
		return lease, nil
placeholder

	if req.ForceNewConn {
		p.recordConnPickDuration(time.Since(pickStartedAt))
		ap.mu.Unlock()
		closeOpenAIWSConns(evicted)
		return nil, errOpenAIWSConnQueueFull
placeholder

acquireAtCapacity:
	target := p.pickLeastBusyConnLocked(ap, req.PreferredConnID, compatibility)
	connPick := time.Since(pickStartedAt)
	p.recordConnPickDuration(connPick)
	if target == nil {
		ap.mu.Unlock()
		closeOpenAIWSConns(evicted)
		return nil, errOpenAIWSConnClosed
placeholder
	if int(target.waiters.Load()) >= p.queueLimitPerConn() {
		ap.mu.Unlock()
		closeOpenAIWSConns(evicted)
		return nil, errOpenAIWSConnQueueFull
placeholder
	target.waiters.Add(1)
	ap.mu.Unlock()
	closeOpenAIWSConns(evicted)
	defer target.waiters.Add(-1)
	waitStart := time.Now()
	p.metrics.acquireQueueWaitTotal.Add(1)

	if err := target.acquire(ctx); err != nil {
		if errors.Is(err, errOpenAIWSConnClosed) && retry < 1 {
			return p.acquire(ctx, req, retry+1)
	placeholder
		return nil, err
placeholder
	if p.shouldHealthCheckConn(target) {
		if err := target.pingWithTimeout(openAIWSConnHealthCheckTO); err != nil {
			target.release()
			target.close()
			p.evictConn(accountID, target.id)
			if retry < 1 {
				return p.acquire(ctx, req, retry+1)
		placeholder
			return nil, err
	placeholder
placeholder

	queueWait := time.Since(waitStart)
	p.metrics.acquireQueueWaitMs.Add(queueWait.Milliseconds())
	lease := &openAIWSConnLease{pool: p, accountID: accountID, conn: target, queueWait: queueWait, connPick: connPick, reused: trueplaceholder
	p.metrics.acquireReuseTotal.Add(1)
	p.recordLastSuccessfulAcquire(accountID, acquireGeneration, req)
	p.ensureTargetIdleAsync(accountID)
	return lease, nil
placeholder

func (p *openAIWSConnPool) recordConnPickDuration(duration time.Duration) {
	if p == nil {
		return
placeholder
	if duration < 0 {
		duration = 0
placeholder
	p.metrics.connPickTotal.Add(1)
	p.metrics.connPickMs.Add(duration.Milliseconds())
placeholder

func (p *openAIWSConnPool) recordLastSuccessfulAcquire(accountID int64, generation uint64, req openAIWSAcquireRequest) {
	if p == nil || accountID <= 0 {
		return
placeholder
	ap, ok := p.getAccountPool(accountID)
	if !ok || ap == nil {
		return
placeholder
	ap.mu.Lock()
	if ap.generation != generation {
		ap.mu.Unlock()
		return
placeholder
	ap.lastAcquire = cloneOpenAIWSAcquireRequestPtr(&req)
	ap.mu.Unlock()
placeholder

func (p *openAIWSConnPool) pickOldestIdleConnLocked(ap *openAIWSAccountPool) *openAIWSConn {
	if ap == nil || len(ap.conns) == 0 {
		return nil
placeholder
	var oldest *openAIWSConn
	for _, conn := range ap.conns {
		if conn == nil || conn.isLeased() || conn.waiters.Load() > 0 || p.isConnPinnedLocked(ap, conn.id) {
			continue
	placeholder
		if oldest == nil || conn.lastUsedAt().Before(oldest.lastUsedAt()) {
			oldest = conn
	placeholder
placeholder
	return oldest
placeholder

func (p *openAIWSConnPool) pickOldestIdleConnWithoutHandshakeCompatibilityOrRoutingAffinityLocked(
	ap *openAIWSAccountPool,
	compatibility openAIWSHandshakeCompatibilityKey,
	routingAffinity string,
) *openAIWSConn {
	if ap == nil || len(ap.conns) == 0 {
		return nil
placeholder
	var oldest *openAIWSConn
	for _, conn := range ap.conns {
		if conn == nil ||
			(conn.matchesHandshakeCompatibility(compatibility) && conn.matchesRoutingAffinity(routingAffinity)) ||
			conn.isLeased() || conn.waiters.Load() > 0 || p.isConnPinnedLocked(ap, conn.id) {
			continue
	placeholder
		if oldest == nil || conn.lastUsedAt().Before(oldest.lastUsedAt()) {
			oldest = conn
	placeholder
placeholder
	return oldest
placeholder

func (p *openAIWSConnPool) getOrCreateAccountPool(accountID int64) *openAIWSAccountPool {
	if p == nil || accountID <= 0 {
		return nil
placeholder
	if existing, ok := p.accounts.Load(accountID); ok {
		if ap, typed := existing.(*openAIWSAccountPool); typed && ap != nil {
			return ap
	placeholder
placeholder
	ap := &openAIWSAccountPool{
		conns:       make(map[string]*openAIWSConn),
		pinnedConns: make(map[string]int),
		changedCh:   make(chan struct{placeholder),
placeholder
	actual, _ := p.accounts.LoadOrStore(accountID, ap)
	if typed, ok := actual.(*openAIWSAccountPool); ok && typed != nil {
		return typed
placeholder
	return ap
placeholder

// ensureAccountPoolLocked 兼容旧调用。
func (p *openAIWSConnPool) ensureAccountPoolLocked(accountID int64) *openAIWSAccountPool {
	return p.getOrCreateAccountPool(accountID)
placeholder

func (p *openAIWSConnPool) getAccountPool(accountID int64) (*openAIWSAccountPool, bool) {
	if p == nil || accountID <= 0 {
		return nil, false
placeholder
	value, ok := p.accounts.Load(accountID)
	if !ok || value == nil {
		return nil, false
placeholder
	ap, typed := value.(*openAIWSAccountPool)
	return ap, typed && ap != nil
placeholder

func (p *openAIWSConnPool) notifyAccountPoolChanged(accountID int64) {
	ap, ok := p.getAccountPool(accountID)
	if !ok || ap == nil {
		return
placeholder
	ap.mu.Lock()
	ap.signalChangedLocked()
	ap.mu.Unlock()
placeholder

func (p *openAIWSConnPool) isConnPinnedLocked(ap *openAIWSAccountPool, connID string) bool {
	if ap == nil || connID == "" || len(ap.pinnedConns) == 0 {
		return false
placeholder
	return ap.pinnedConns[connID] > 0
placeholder

func (p *openAIWSConnPool) cleanupAccountLocked(ap *openAIWSAccountPool, now time.Time, maxConns int) []*openAIWSConn {
	if ap == nil {
		return nil
placeholder
	maxAge := p.maxConnAge()

	evicted := make([]*openAIWSConn, 0)
	for id, conn := range ap.conns {
		if conn == nil {
			delete(ap.conns, id)
			if len(ap.pinnedConns) > 0 {
				delete(ap.pinnedConns, id)
		placeholder
			continue
	placeholder
		select {
		case <-conn.closedCh:
			delete(ap.conns, id)
			if len(ap.pinnedConns) > 0 {
				delete(ap.pinnedConns, id)
		placeholder
			evicted = append(evicted, conn)
			continue
		default:
	placeholder
		if p.isConnPinnedLocked(ap, id) {
			continue
	placeholder
		if maxAge > 0 && !conn.isLeased() && conn.age(now) > maxAge {
			delete(ap.conns, id)
			if len(ap.pinnedConns) > 0 {
				delete(ap.pinnedConns, id)
		placeholder
			evicted = append(evicted, conn)
	placeholder
placeholder

	if maxConns <= 0 {
		maxConns = p.maxConnsHardCap()
placeholder
	maxIdle := p.maxIdlePerAccount()
	if maxIdle < 0 || maxIdle > maxConns {
		maxIdle = maxConns
placeholder
	if maxIdle >= 0 && len(ap.conns) > maxIdle {
		idleConns := make([]*openAIWSConn, 0, len(ap.conns))
		for id, conn := range ap.conns {
			if conn == nil {
				delete(ap.conns, id)
				if len(ap.pinnedConns) > 0 {
					delete(ap.pinnedConns, id)
			placeholder
				continue
		placeholder
			// 有等待者的连接不能在清理阶段被淘汰，否则等待中的 acquire 会收到 closed 错误。
			if conn.isLeased() || conn.waiters.Load() > 0 || p.isConnPinnedLocked(ap, conn.id) {
				continue
		placeholder
			idleConns = append(idleConns, conn)
	placeholder
		sort.SliceStable(idleConns, func(i, j int) bool {
			return idleConns[i].lastUsedAt().Before(idleConns[j].lastUsedAt())
	placeholder)
		redundant := len(ap.conns) - maxIdle
		if redundant > len(idleConns) {
			redundant = len(idleConns)
	placeholder
		for i := 0; i < redundant; i++ {
			conn := idleConns[i]
			delete(ap.conns, conn.id)
			if len(ap.pinnedConns) > 0 {
				delete(ap.pinnedConns, conn.id)
		placeholder
			evicted = append(evicted, conn)
	placeholder
		if redundant > 0 {
			p.metrics.scaleDownTotal.Add(int64(redundant))
	placeholder
placeholder
	if len(evicted) > 0 {
		ap.signalChangedLocked()
placeholder

	return evicted
placeholder

func (p *openAIWSConnPool) pickLeastBusyConnLocked(
	ap *openAIWSAccountPool,
	preferredConnID string,
	compatibility openAIWSHandshakeCompatibilityKey,
) *openAIWSConn {
	if ap == nil || len(ap.conns) == 0 {
		return nil
placeholder
	preferredConnID = stringsTrim(preferredConnID)
	if preferredConnID != "" {
		if conn, ok := ap.conns[preferredConnID]; ok && conn.matchesHandshakeCompatibility(compatibility) {
			return conn
	placeholder
placeholder
	var best *openAIWSConn
	var bestWaiters int32
	var bestLastUsed time.Time
	for _, conn := range ap.conns {
		if conn == nil || !conn.matchesHandshakeCompatibility(compatibility) {
			continue
	placeholder
		waiters := conn.waiters.Load()
		lastUsed := conn.lastUsedAt()
		if best == nil ||
			waiters < bestWaiters ||
			(waiters == bestWaiters && lastUsed.Before(bestLastUsed)) {
			best = conn
			bestWaiters = waiters
			bestLastUsed = lastUsed
	placeholder
placeholder
	return best
placeholder

func (p *openAIWSConnPool) pickLeastBusyConnWithRoutingAffinityLocked(
	ap *openAIWSAccountPool,
	compatibility openAIWSHandshakeCompatibilityKey,
	routingAffinity string,
) *openAIWSConn {
	if ap == nil || len(ap.conns) == 0 {
		return nil
placeholder
	var best *openAIWSConn
	var bestWaiters int32
	var bestLastUsed time.Time
	for _, conn := range ap.conns {
		if conn == nil ||
			!conn.matchesHandshakeCompatibility(compatibility) ||
			!conn.matchesRoutingAffinity(routingAffinity) {
			continue
	placeholder
		waiters := conn.waiters.Load()
		lastUsed := conn.lastUsedAt()
		if best == nil ||
			waiters < bestWaiters ||
			(waiters == bestWaiters && lastUsed.Before(bestLastUsed)) {
			best = conn
			bestWaiters = waiters
			bestLastUsed = lastUsed
	placeholder
placeholder
	return best
placeholder

func accountPoolLoadLocked(ap *openAIWSAccountPool) (inflight int, waiters int) {
	if ap == nil {
		return 0, 0
placeholder
	for _, conn := range ap.conns {
		if conn == nil {
			continue
	placeholder
		if conn.isLeased() {
			inflight++
	placeholder
		waiters += int(conn.waiters.Load())
placeholder
	return inflight, waiters
placeholder

// AccountPoolLoad 返回指定账号连接池的并发与排队快照。
func (p *openAIWSConnPool) AccountPoolLoad(accountID int64) (inflight int, waiters int, conns int) {
	if p == nil || accountID <= 0 {
		return 0, 0, 0
placeholder
	ap, ok := p.getAccountPool(accountID)
	if !ok || ap == nil {
		return 0, 0, 0
placeholder
	ap.mu.Lock()
	defer ap.mu.Unlock()
	inflight, waiters = accountPoolLoadLocked(ap)
	return inflight, waiters, len(ap.conns)
placeholder

func (p *openAIWSConnPool) ensureTargetIdleAsync(accountID int64) {
	if p == nil || accountID <= 0 {
		return
placeholder

	var req openAIWSAcquireRequest
	generation := uint64(0)
	need := 0
	ap, ok := p.getAccountPool(accountID)
	if !ok || ap == nil {
		return
placeholder
	ap.mu.Lock()
	defer ap.mu.Unlock()
	if ap.lastAcquire == nil {
		return
placeholder
	if ap.prewarmActive {
		return
placeholder
	now := time.Now()
	if !ap.prewarmUntil.IsZero() && now.Before(ap.prewarmUntil) {
		return
placeholder
	if p.shouldSuppressPrewarmLocked(ap, now) {
		return
placeholder
	effectiveMaxConns := p.maxConnsHardCap()
	if ap.lastAcquire != nil && ap.lastAcquire.Account != nil {
		effectiveMaxConns = p.effectiveMaxConnsByAccount(ap.lastAcquire.Account)
placeholder
	target := p.targetConnCountLocked(ap, effectiveMaxConns)
	current := len(ap.conns) + ap.creating
	if current >= target {
		return
placeholder
	need = target - current
	if need <= 0 {
		return
placeholder
	req = cloneOpenAIWSAcquireRequest(*ap.lastAcquire)
	generation = ap.generation
	ap.prewarmActive = true
	if cooldown := p.prewarmCooldown(); cooldown > 0 {
		ap.prewarmUntil = now.Add(cooldown)
placeholder
	ap.creating += need
	p.metrics.scaleUpTotal.Add(int64(need))

	go p.prewarmConns(accountID, req, need, generation)
placeholder

func (p *openAIWSConnPool) targetConnCountLocked(ap *openAIWSAccountPool, maxConns int) int {
	if ap == nil {
		return 0
placeholder

	if maxConns <= 0 {
		return 0
placeholder

	minIdle := p.minIdlePerAccount()
	if minIdle < 0 {
		minIdle = 0
placeholder
	if minIdle > maxConns {
		minIdle = maxConns
placeholder

	inflight, waiters := accountPoolLoadLocked(ap)
	utilization := p.targetUtilization()
	demand := inflight + waiters
	if demand <= 0 {
		return minIdle
placeholder

	target := 1
	if demand > 1 {
		target = int(math.Ceil(float64(demand) / utilization))
placeholder
	if waiters > 0 && target < len(ap.conns)+1 {
		target = len(ap.conns) + 1
placeholder
	if target < minIdle {
		target = minIdle
placeholder
	if target > maxConns {
		target = maxConns
placeholder
	return target
placeholder

func (p *openAIWSConnPool) prewarmConns(accountID int64, req openAIWSAcquireRequest, total int, generations ...uint64) {
	generation := uint64(0)
	if len(generations) > 0 {
		generation = generations[0]
placeholder
	staleTarget := false
	defer func() {
		if ap, ok := p.getAccountPool(accountID); ok && ap != nil {
			ap.mu.Lock()
			ap.prewarmActive = false
			ap.signalChangedLocked()
			ap.mu.Unlock()
	placeholder
		if staleTarget {
			// A newer acquire arrived while the old dial was in flight. Re-run
			// target selection only after clearing prewarmActive so the latest
			// beta/hint target can fill the idle budget.
			p.ensureTargetIdleAsync(accountID)
	placeholder
placeholder()

	for i := 0; i < total; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), p.dialTimeout()+openAIWSConnPrewarmExtraDelay)
		conn, err := p.dialConn(ctx, req)
		cancel()

		ap, ok := p.getAccountPool(accountID)
		if !ok || ap == nil {
			if conn != nil {
				conn.close()
		placeholder
			return
	placeholder
		ap.mu.Lock()
		if ap.creating > 0 {
			ap.creating--
	placeholder
		if err != nil {
			ap.prewarmFails++
			ap.prewarmFailAt = time.Now()
			ap.signalChangedLocked()
			ap.mu.Unlock()
			continue
	placeholder
		if ap.generation != generation || ap.lastAcquire == nil {
			ap.mu.Unlock()
			conn.close()
			continue
	placeholder
		if !sameOpenAIWSPrewarmTarget(req, *ap.lastAcquire) {
			staleTarget = true
			ap.signalChangedLocked()
			ap.mu.Unlock()
			conn.close()
			continue
	placeholder
		if len(ap.conns) >= p.effectiveMaxConnsByAccount(req.Account) {
			ap.signalChangedLocked()
			ap.mu.Unlock()
			conn.close()
			continue
	placeholder
		ap.conns[conn.id] = conn
		ap.prewarmFails = 0
		ap.prewarmFailAt = time.Time{placeholder
		ap.signalChangedLocked()
		ap.mu.Unlock()
placeholder
placeholder

// ClearAccount closes all pooled connections and discards delayed prewarm
// state for one account. The generation guard prevents an in-flight prewarm
// started before credential recovery from re-entering the pool afterwards.
func (p *openAIWSConnPool) ClearAccount(accountID int64) {
	if p == nil || accountID <= 0 {
		return
placeholder
	ap, ok := p.getAccountPool(accountID)
	if !ok || ap == nil {
		return
placeholder
	ap.mu.Lock()
	ap.generation++
	conns := make([]*openAIWSConn, 0, len(ap.conns))
	for id, conn := range ap.conns {
		delete(ap.conns, id)
		delete(ap.pinnedConns, id)
		if conn != nil {
			conns = append(conns, conn)
	placeholder
placeholder
	ap.lastAcquire = nil
	ap.prewarmUntil = time.Time{placeholder
	ap.prewarmFails = 0
	ap.prewarmFailAt = time.Time{placeholder
	ap.signalChangedLocked()
	ap.mu.Unlock()
	closeOpenAIWSConns(conns)
placeholder

func (p *openAIWSConnPool) evictConn(accountID int64, connID string) {
	if p == nil || accountID <= 0 || stringsTrim(connID) == "" {
		return
placeholder
	var conn *openAIWSConn
	ap, ok := p.getAccountPool(accountID)
	if ok && ap != nil {
		ap.mu.Lock()
		if c, exists := ap.conns[connID]; exists {
			conn = c
			delete(ap.conns, connID)
			if len(ap.pinnedConns) > 0 {
				delete(ap.pinnedConns, connID)
		placeholder
			ap.signalChangedLocked()
	placeholder
		ap.mu.Unlock()
placeholder
	if conn != nil {
		conn.close()
placeholder
placeholder

func (p *openAIWSConnPool) PinConn(accountID int64, connID string) bool {
	if p == nil || accountID <= 0 {
		return false
placeholder
	connID = stringsTrim(connID)
	if connID == "" {
		return false
placeholder
	ap, ok := p.getAccountPool(accountID)
	if !ok || ap == nil {
		return false
placeholder
	ap.mu.Lock()
	defer ap.mu.Unlock()
	if _, exists := ap.conns[connID]; !exists {
		return false
placeholder
	if ap.pinnedConns == nil {
		ap.pinnedConns = make(map[string]int)
placeholder
	ap.pinnedConns[connID]++
	return true
placeholder

func (p *openAIWSConnPool) UnpinConn(accountID int64, connID string) {
	if p == nil || accountID <= 0 {
		return
placeholder
	connID = stringsTrim(connID)
	if connID == "" {
		return
placeholder
	ap, ok := p.getAccountPool(accountID)
	if !ok || ap == nil {
		return
placeholder
	ap.mu.Lock()
	defer ap.mu.Unlock()
	if len(ap.pinnedConns) == 0 {
		return
placeholder
	count := ap.pinnedConns[connID]
	if count <= 1 {
		delete(ap.pinnedConns, connID)
		ap.signalChangedLocked()
		return
placeholder
	ap.pinnedConns[connID] = count - 1
	ap.signalChangedLocked()
placeholder

func (p *openAIWSConnPool) dialConn(ctx context.Context, req openAIWSAcquireRequest) (*openAIWSConn, error) {
	if p == nil || p.clientDialer == nil {
		return nil, errors.New("openai ws client dialer is nil")
placeholder
	headers := cloneHeader(req.Headers)
	var err error
	if req.HeadersFactory != nil {
		headers, err = req.HeadersFactory(ctx, headers)
		if err != nil {
			return nil, err
	placeholder
placeholder
	conn, status, handshakeHeaders, err := p.clientDialer.Dial(ctx, req.WSURL, headers, req.ProxyURL)
	if err != nil {
		var handshakeErr *openAIWSHandshakeError
		var responseBody []byte
		if errors.As(err, &handshakeErr) && handshakeErr != nil {
			responseBody = append([]byte(nil), handshakeErr.Body...)
	placeholder
		return nil, &openAIWSDialError{
			StatusCode:      status,
			ResponseHeaders: cloneHeader(handshakeHeaders),
			ResponseBody:    responseBody,
			Err:             err,
	placeholder
placeholder
	if conn == nil {
		return nil, &openAIWSDialError{
			StatusCode:      status,
			ResponseHeaders: cloneHeader(handshakeHeaders),
			Err:             errors.New("openai ws dialer returned nil connection"),
	placeholder
placeholder
	id := p.nextConnID(req.Account.ID)
	pooledConn := newOpenAIWSConn(id, req.Account.ID, conn, handshakeHeaders)
	pooledConn.handshakeCompatibility = normalizeOpenAIWSHandshakeCompatibility(req.Headers)
	pooledConn.routingAffinity = normalizeOpenAIWSRoutingAffinity(req.Headers)
	return pooledConn, nil
placeholder

func (p *openAIWSConnPool) nextConnID(accountID int64) string {
	seq := p.seq.Add(1)
	buf := make([]byte, 0, 32)
	buf = append(buf, "oa_ws_"...)
	buf = strconv.AppendInt(buf, accountID, 10)
	buf = append(buf, '_')
	buf = strconv.AppendUint(buf, seq, 10)
	return string(buf)
placeholder

func (p *openAIWSConnPool) shouldHealthCheckConn(conn *openAIWSConn) bool {
	if conn == nil || !conn.supportsIdlePingWithoutReader() {
		return false
placeholder
	return conn.idleDuration(time.Now()) >= openAIWSConnHealthCheckIdle
placeholder

func (p *openAIWSConnPool) maxConnsHardCap() int {
	if p != nil && p.cfg != nil && p.cfg.Gateway.OpenAIWS.MaxConnsPerAccount > 0 {
		return p.cfg.Gateway.OpenAIWS.MaxConnsPerAccount
placeholder
	return 8
placeholder

func (p *openAIWSConnPool) dynamicMaxConnsEnabled() bool {
	if p != nil && p.cfg != nil {
		return p.cfg.Gateway.OpenAIWS.DynamicMaxConnsByAccountConcurrencyEnabled
placeholder
	return false
placeholder

func (p *openAIWSConnPool) modeRouterV2Enabled() bool {
	if p != nil && p.cfg != nil {
		return p.cfg.Gateway.OpenAIWS.ModeRouterV2Enabled
placeholder
	return false
placeholder

func (p *openAIWSConnPool) maxConnsFactorByAccount(account *Account) float64 {
	if p == nil || p.cfg == nil || account == nil {
		return 1.0
placeholder
	switch account.Type {
	case AccountTypeOAuth:
		if p.cfg.Gateway.OpenAIWS.OAuthMaxConnsFactor > 0 {
			return p.cfg.Gateway.OpenAIWS.OAuthMaxConnsFactor
	placeholder
	case AccountTypeAPIKey:
		if p.cfg.Gateway.OpenAIWS.APIKeyMaxConnsFactor > 0 {
			return p.cfg.Gateway.OpenAIWS.APIKeyMaxConnsFactor
	placeholder
placeholder
	return 1.0
placeholder

func (p *openAIWSConnPool) effectiveMaxConnsByAccount(account *Account) int {
	hardCap := p.maxConnsHardCap()
	if hardCap <= 0 {
		return 0
placeholder
	if p.modeRouterV2Enabled() {
		if account == nil {
			return hardCap
	placeholder
		if account.Concurrency <= 0 {
			return 0
	placeholder
		return min(account.Concurrency, hardCap)
placeholder
	if account == nil || !p.dynamicMaxConnsEnabled() {
		return hardCap
placeholder
	if account.Concurrency <= 0 {
		// 0/-1 等“无限制”并发场景下，仍由全局硬上限兜底。
		return hardCap
placeholder
	factor := p.maxConnsFactorByAccount(account)
	if factor <= 0 {
		factor = 1.0
placeholder
	effective := int(math.Ceil(float64(account.Concurrency) * factor))
	if effective < 1 {
		effective = 1
placeholder
	if effective > hardCap {
		effective = hardCap
placeholder
	return effective
placeholder

func (p *openAIWSConnPool) minIdlePerAccount() int {
	if p != nil && p.cfg != nil && p.cfg.Gateway.OpenAIWS.MinIdlePerAccount >= 0 {
		return p.cfg.Gateway.OpenAIWS.MinIdlePerAccount
placeholder
	return 0
placeholder

func (p *openAIWSConnPool) maxIdlePerAccount() int {
	if p != nil && p.cfg != nil && p.cfg.Gateway.OpenAIWS.MaxIdlePerAccount >= 0 {
		return p.cfg.Gateway.OpenAIWS.MaxIdlePerAccount
placeholder
	return 4
placeholder

func (p *openAIWSConnPool) maxConnAge() time.Duration {
	return openAIWSConnMaxAge
placeholder

func (p *openAIWSConnPool) queueLimitPerConn() int {
	if p != nil && p.cfg != nil && p.cfg.Gateway.OpenAIWS.QueueLimitPerConn > 0 {
		return p.cfg.Gateway.OpenAIWS.QueueLimitPerConn
placeholder
	return 256
placeholder

func (p *openAIWSConnPool) targetUtilization() float64 {
	if p != nil && p.cfg != nil {
		ratio := p.cfg.Gateway.OpenAIWS.PoolTargetUtilization
		if ratio > 0 && ratio <= 1 {
			return ratio
	placeholder
placeholder
	return 0.7
placeholder

func (p *openAIWSConnPool) prewarmCooldown() time.Duration {
	if p != nil && p.cfg != nil && p.cfg.Gateway.OpenAIWS.PrewarmCooldownMS > 0 {
		return time.Duration(p.cfg.Gateway.OpenAIWS.PrewarmCooldownMS) * time.Millisecond
placeholder
	return 0
placeholder

func (p *openAIWSConnPool) shouldSuppressPrewarmLocked(ap *openAIWSAccountPool, now time.Time) bool {
	if ap == nil {
		return true
placeholder
	if ap.prewarmFails <= 0 {
		return false
placeholder
	if ap.prewarmFailAt.IsZero() {
		ap.prewarmFails = 0
		return false
placeholder
	if now.Sub(ap.prewarmFailAt) > openAIWSPrewarmFailureWindow {
		ap.prewarmFails = 0
		ap.prewarmFailAt = time.Time{placeholder
		return false
placeholder
	return ap.prewarmFails >= openAIWSPrewarmFailureSuppress
placeholder

func (p *openAIWSConnPool) dialTimeout() time.Duration {
	if p != nil && p.cfg != nil && p.cfg.Gateway.OpenAIWS.DialTimeoutSeconds > 0 {
		return time.Duration(p.cfg.Gateway.OpenAIWS.DialTimeoutSeconds) * time.Second
placeholder
	return 10 * time.Second
placeholder

func cloneOpenAIWSAcquireRequest(req openAIWSAcquireRequest) openAIWSAcquireRequest {
	copied := req
	copied.Headers = cloneHeader(req.Headers)
	copied.WSURL = stringsTrim(req.WSURL)
	copied.ProxyURL = stringsTrim(req.ProxyURL)
	copied.PreferredConnID = stringsTrim(req.PreferredConnID)
	return copied
placeholder

func cloneOpenAIWSAcquireRequestPtr(req *openAIWSAcquireRequest) *openAIWSAcquireRequest {
	if req == nil {
		return nil
placeholder
	copied := cloneOpenAIWSAcquireRequest(*req)
	return &copied
placeholder

func sameOpenAIWSPrewarmTarget(a, b openAIWSAcquireRequest) bool {
	return stringsTrim(a.WSURL) == stringsTrim(b.WSURL) &&
		stringsTrim(a.ProxyURL) == stringsTrim(b.ProxyURL) &&
		normalizeOpenAIWSHandshakeCompatibility(a.Headers) == normalizeOpenAIWSHandshakeCompatibility(b.Headers)
placeholder

func normalizeOpenAIWSBetaFeatures(headers http.Header) string {
	features := make(map[string]struct{placeholder)
	for name, values := range headers {
		if !strings.EqualFold(strings.TrimSpace(name), "x-codex-beta-features") {
			continue
	placeholder
		for _, value := range values {
			for _, feature := range strings.Split(value, ",") {
				if feature = strings.TrimSpace(feature); feature != "" {
					features[feature] = struct{placeholder{placeholder
			placeholder
		placeholder
	placeholder
placeholder
	if len(features) == 0 {
		return ""
placeholder
	normalized := make([]string, 0, len(features))
	for feature := range features {
		normalized = append(normalized, feature)
placeholder
	sort.Strings(normalized)
	return strings.Join(normalized, ",")
placeholder

func normalizeOpenAIWSHandshakeCompatibility(headers http.Header) openAIWSHandshakeCompatibilityKey {
	return openAIWSHandshakeCompatibilityKey{
		betaFeatures: normalizeOpenAIWSBetaFeatures(headers),
placeholder
placeholder

func normalizeOpenAIWSRoutingAffinity(headers http.Header) string {
	canonicalName := http.CanonicalHeaderKey(openAICodexRoutingHintHeader)
	if values, ok := headers[canonicalName]; ok {
		for _, value := range values {
			if value = strings.TrimSpace(value); value != "" {
				return value
		placeholder
	placeholder
placeholder

	variantNames := make([]string, 0)
	for name := range headers {
		if name != canonicalName && strings.EqualFold(strings.TrimSpace(name), openAICodexRoutingHintHeader) {
			variantNames = append(variantNames, name)
	placeholder
placeholder
	sort.Strings(variantNames)
	for _, name := range variantNames {
		for _, value := range headers[name] {
			if value = strings.TrimSpace(value); value != "" {
				return value
		placeholder
	placeholder
placeholder
	return ""
placeholder

func cloneHeader(src http.Header) http.Header {
	if src == nil {
		return nil
placeholder
	dst := make(http.Header, len(src))
	for k, vals := range src {
		if len(vals) == 0 {
			dst[k] = nil
			continue
	placeholder
		copied := make([]string, len(vals))
		copy(copied, vals)
		dst[k] = copied
placeholder
	return dst
placeholder

func closeOpenAIWSConns(conns []*openAIWSConn) {
	if len(conns) == 0 {
		return
placeholder
	for _, conn := range conns {
		if conn == nil {
			continue
	placeholder
		conn.close()
placeholder
placeholder

func stringsTrim(value string) string {
	return strings.TrimSpace(value)
placeholder
