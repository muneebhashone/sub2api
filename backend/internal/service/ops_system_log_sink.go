package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/util/logredact"
)

type OpsSystemLogSinkHealth struct {
	QueueDepth      int64  `json:"queue_depth"`
	QueueCapacity   int64  `json:"queue_capacity"`
	DroppedCount    uint64 `json:"dropped_count"`
	WriteFailed     uint64 `json:"write_failed_count"`
	WrittenCount    uint64 `json:"written_count"`
	AvgWriteDelayMs uint64 `json:"avg_write_delay_ms"`
	LastError       string `json:"last_error"`
placeholder

type OpsSystemLogSink struct {
	opsRepo OpsRepository
	host    string

	queue chan *logger.LogEvent

	batchSize     int
	flushInterval time.Duration

	// 连续写入失败后的退避参数。构造后只读，测试可在 Start 前覆盖。
	flushBackoff    time.Duration
	flushBackoffMax time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	droppedCount uint64
	writeFailed  uint64
	writtenCount uint64
	totalDelayNs uint64

	lastError atomic.Value
placeholder

const maxSystemLogHostLength = 255

const (
	// 首次写入失败后暂停落库的时长，之后逐次翻倍到上限。
	defaultOpsSystemLogFlushBackoff = 2 * time.Second
	// 退避上限。日志是尽力而为的观测数据，不值得为它无限期占用连接池。
	defaultOpsSystemLogFlushBackoffMax = 60 * time.Second
)

func NewOpsSystemLogSink(opsRepo OpsRepository) *OpsSystemLogSink {
	ctx, cancel := context.WithCancel(context.Background())
	rawHost, err := os.Hostname()
	s := &OpsSystemLogSink{
		opsRepo:         opsRepo,
		host:            normalizeSystemLogHost(rawHost, err),
		queue:           make(chan *logger.LogEvent, 5000),
		batchSize:       200,
		flushInterval:   time.Second,
		flushBackoff:    defaultOpsSystemLogFlushBackoff,
		flushBackoffMax: defaultOpsSystemLogFlushBackoffMax,
		ctx:             ctx,
		cancel:          cancel,
placeholder
	s.lastError.Store("")
	return s
placeholder

// flushBackoffFor 返回第 failures 次连续失败后的退避时长（指数退避，封顶）。
func (s *OpsSystemLogSink) flushBackoffFor(failures int) time.Duration {
	base := s.flushBackoff
	if base <= 0 {
		base = defaultOpsSystemLogFlushBackoff
placeholder
	maxBackoff := s.flushBackoffMax
	if maxBackoff <= 0 {
		maxBackoff = defaultOpsSystemLogFlushBackoffMax
placeholder
	if maxBackoff < base {
		maxBackoff = base
placeholder
	backoff := base
	for i := 1; i < failures && backoff < maxBackoff; i++ {
		backoff *= 2
placeholder
	if backoff > maxBackoff {
		backoff = maxBackoff
placeholder
	return backoff
placeholder

func normalizeSystemLogHost(host string, err error) string {
	host = strings.TrimSpace(host)
	if err != nil || host == "" {
		return "unknown"
placeholder
	runes := []rune(host)
	if len(runes) > maxSystemLogHostLength {
		return string(runes[:maxSystemLogHostLength])
placeholder
	return host
placeholder

func (s *OpsSystemLogSink) Start() {
	if s == nil || s.opsRepo == nil {
		return
placeholder
	s.wg.Add(1)
	go s.run()
placeholder

func (s *OpsSystemLogSink) Stop() {
	if s == nil {
		return
placeholder
	s.cancel()
	s.wg.Wait()
placeholder

func (s *OpsSystemLogSink) WriteLogEvent(event *logger.LogEvent) {
	if s == nil || event == nil || !s.shouldIndex(event) {
		return
placeholder
	if s.ctx != nil {
		select {
		case <-s.ctx.Done():
			return
		default:
	placeholder
placeholder

	select {
	case s.queue <- event:
	default:
		atomic.AddUint64(&s.droppedCount, 1)
placeholder
placeholder

func (s *OpsSystemLogSink) shouldIndex(event *logger.LogEvent) bool {
	if event != nil && event.Fields != nil {
		if skip, _ := event.Fields[logger.OpsSystemLogSkipField].(bool); skip {
			return false
	placeholder
placeholder
	level := strings.ToLower(strings.TrimSpace(event.Level))
	switch level {
	case "warn", "warning", "error", "fatal", "panic", "dpanic":
		return true
placeholder

	component := strings.ToLower(strings.TrimSpace(event.Component))
	// zap 的 LoggerName 往往为空或不等于业务组件名；业务组件名通常以字段 component 透传。
	if event.Fields != nil {
		if fc := strings.ToLower(strings.TrimSpace(asString(event.Fields["component"]))); fc != "" {
			component = fc
	placeholder
placeholder
	if strings.Contains(component, "http.access") {
		return true
placeholder
	if strings.Contains(component, "audit") {
		return true
placeholder
	return false
placeholder

func (s *OpsSystemLogSink) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.flushInterval)
	defer ticker.Stop()

	batch := make([]*logger.LogEvent, 0, s.batchSize)
	// 仅在本 goroutine 内读写，无需加锁。
	failures := 0
	var suppressedUntil time.Time
	flush := func(baseCtx context.Context) {
		if len(batch) == 0 {
			return
	placeholder
		now := time.Now()
		if now.Before(suppressedUntil) {
			// 退避窗口内直接丢弃本批：日志是尽力而为的观测数据，继续攒批只会把
			// 压力转移到内存，而每次重试都会再占用并取消一条池内连接。
			atomic.AddUint64(&s.droppedCount, uint64(len(batch)))
			batch = batch[:0]
			return
	placeholder
		started := time.Now()
		inserted, err := s.flushBatch(baseCtx, batch)
		delay := time.Since(started)
		if err != nil {
			failures++
			backoff := s.flushBackoffFor(failures)
			suppressedUntil = time.Now().Add(backoff)
			atomic.AddUint64(&s.writeFailed, uint64(len(batch)))
			s.lastError.Store(err.Error())
			// 每个退避窗口至多一条，避免数据库故障期间刷屏。
			_, _ = fmt.Fprintf(os.Stderr, "time=%s level=WARN msg=\"ops system log sink flush failed\" err=%v batch=%d failures=%d backoff=%s\n",
				time.Now().Format(time.RFC3339Nano), err, len(batch), failures, backoff,
			)
	placeholder else {
			failures = 0
			suppressedUntil = time.Time{placeholder
			atomic.AddUint64(&s.writtenCount, uint64(inserted))
			atomic.AddUint64(&s.totalDelayNs, uint64(delay.Nanoseconds()))
			s.lastError.Store("")
	placeholder
		batch = batch[:0]
placeholder
	drainAndFlush := func() {
		for {
			select {
			case item := <-s.queue:
				if item == nil {
					continue
			placeholder
				batch = append(batch, item)
				if len(batch) >= s.batchSize {
					flush(context.Background())
			placeholder
			default:
				flush(context.Background())
				return
		placeholder
	placeholder
placeholder

	for {
		select {
		case <-s.ctx.Done():
			drainAndFlush()
			return
		case item := <-s.queue:
			if item == nil {
				continue
		placeholder
			batch = append(batch, item)
			if len(batch) >= s.batchSize {
				flush(s.ctx)
		placeholder
		case <-ticker.C:
			flush(s.ctx)
	placeholder
placeholder
placeholder

func (s *OpsSystemLogSink) flushBatch(baseCtx context.Context, batch []*logger.LogEvent) (int, error) {
	inputs := make([]*OpsInsertSystemLogInput, 0, len(batch))
	for _, event := range batch {
		if event == nil {
			continue
	placeholder
		createdAt := event.Time.UTC()
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
	placeholder

		fields := copyMap(event.Fields)
		requestID := asString(fields["request_id"])
		clientRequestID := asString(fields["client_request_id"])
		platform := asString(fields["platform"])
		model := asString(fields["model"])
		component := strings.TrimSpace(event.Component)
		if fieldComponent := asString(fields["component"]); fieldComponent != "" {
			component = fieldComponent
	placeholder
		if component == "" {
			component = "app"
	placeholder

		userID := asInt64Ptr(fields["user_id"])
		apiKeyID := asInt64Ptr(fields["api_key_id"])
		accountID := asInt64Ptr(fields["account_id"])

		// 统一脱敏后写入索引。
		message := logredact.RedactText(strings.TrimSpace(event.Message))
		redactedExtra := logredact.RedactMap(fields)
		extraJSONBytes, _ := json.Marshal(redactedExtra)
		extraJSON := string(extraJSONBytes)
		if strings.TrimSpace(extraJSON) == "" {
			extraJSON = "{placeholder"
	placeholder

		inputs = append(inputs, &OpsInsertSystemLogInput{
			CreatedAt:       createdAt,
			Host:            s.host,
			Level:           strings.ToLower(strings.TrimSpace(event.Level)),
			Component:       component,
			Message:         message,
			RequestID:       requestID,
			ClientRequestID: clientRequestID,
			UserID:          userID,
			APIKeyID:        apiKeyID,
			AccountID:       accountID,
			Platform:        platform,
			Model:           model,
			ExtraJSON:       extraJSON,
	placeholder)
placeholder

	if len(inputs) == 0 {
		return 0, nil
placeholder
	if baseCtx == nil || baseCtx.Err() != nil {
		baseCtx = context.Background()
placeholder
	ctx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
	defer cancel()
	inserted, err := s.opsRepo.BatchInsertSystemLogs(ctx, inputs)
	if err != nil {
		return 0, err
placeholder
	return int(inserted), nil
placeholder

func (s *OpsSystemLogSink) Health() OpsSystemLogSinkHealth {
	if s == nil {
		return OpsSystemLogSinkHealth{placeholder
placeholder
	written := atomic.LoadUint64(&s.writtenCount)
	totalDelay := atomic.LoadUint64(&s.totalDelayNs)
	var avgDelay uint64
	if written > 0 {
		avgDelay = (totalDelay / written) / uint64(time.Millisecond)
placeholder

	lastErr, _ := s.lastError.Load().(string)
	return OpsSystemLogSinkHealth{
		QueueDepth:      int64(len(s.queue)),
		QueueCapacity:   int64(cap(s.queue)),
		DroppedCount:    atomic.LoadUint64(&s.droppedCount),
		WriteFailed:     atomic.LoadUint64(&s.writeFailed),
		WrittenCount:    written,
		AvgWriteDelayMs: avgDelay,
		LastError:       strings.TrimSpace(lastErr),
placeholder
placeholder

func copyMap(in map[string]any) map[string]any {
	if len(in) == 0 {
		return map[string]any{placeholder
placeholder
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
placeholder
	return out
placeholder

func asString(v any) string {
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case fmt.Stringer:
		return strings.TrimSpace(t.String())
	default:
		return ""
placeholder
placeholder

func asInt64Ptr(v any) *int64 {
	switch t := v.(type) {
	case int:
		n := int64(t)
		if n <= 0 {
			return nil
	placeholder
		return &n
	case int64:
		n := t
		if n <= 0 {
			return nil
	placeholder
		return &n
	case float64:
		n := int64(t)
		if n <= 0 {
			return nil
	placeholder
		return &n
	case json.Number:
		if n, err := t.Int64(); err == nil {
			if n <= 0 {
				return nil
		placeholder
			return &n
	placeholder
	case string:
		raw := strings.TrimSpace(t)
		if raw == "" {
			return nil
	placeholder
		if n, err := strconv.ParseInt(raw, 10, 64); err == nil {
			if n <= 0 {
				return nil
		placeholder
			return &n
	placeholder
placeholder
	return nil
placeholder
