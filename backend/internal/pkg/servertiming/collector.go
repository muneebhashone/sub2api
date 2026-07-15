package servertiming

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	HeaderName       = "Server-Timing"
	AdminUIHeader    = "X-Admin-UI-Request"
	UserUIHeader     = "X-User-UI-Request"
	MetricDatabase   = "db"
	MetricRedis      = "redis"
	dependencyPrefix = "dep_"

	maxMetricNameLength = 48
	maxIntervals        = 2048
	maxHeaderLength     = 4096
)

type contextKey struct{placeholder

type interval struct {
	start time.Time
	end   time.Time
placeholder

type metric struct {
	count     int64
	intervals []interval
placeholder

// Collector stores request-scoped timing samples. It is safe for concurrent use.
type Collector struct {
	startedAt time.Time

	mu          sync.Mutex
	metrics     map[string]*metric
	cacheStatus string
placeholder

// New creates a collector whose total duration starts at startedAt.
func New(startedAt time.Time) *Collector {
	if startedAt.IsZero() {
		startedAt = time.Now()
placeholder
	return &Collector{
		startedAt: startedAt,
		metrics:   make(map[string]*metric),
placeholder
placeholder

// WithCollector attaches a collector to a context.
func WithCollector(ctx context.Context, collector *Collector) context.Context {
	if ctx == nil {
		ctx = context.Background()
placeholder
	if collector == nil {
		return ctx
placeholder
	return context.WithValue(ctx, contextKey{placeholder, collector)
placeholder

// FromContext returns the request timing collector, when one is active.
func FromContext(ctx context.Context) (*Collector, bool) {
	if ctx == nil {
		return nil, false
placeholder
	collector, ok := ctx.Value(contextKey{placeholder).(*Collector)
	return collector, ok && collector != nil
placeholder

// Active reports whether timing collection is enabled for this request.
func Active(ctx context.Context) bool {
	_, ok := FromContext(ctx)
	return ok
placeholder

// Record adds a completed interval and operation count to a metric.
func Record(ctx context.Context, name string, startedAt, endedAt time.Time, count int) {
	collector, ok := FromContext(ctx)
	if !ok {
		return
placeholder
	collector.Record(name, startedAt, endedAt, count)
placeholder

// RecordInterval adds timing without incrementing the operation count. It is
// useful when one logical operation has multiple blocking driver calls.
func RecordInterval(ctx context.Context, name string, startedAt, endedAt time.Time) {
	collector, ok := FromContext(ctx)
	if !ok {
		return
placeholder
	collector.record(name, startedAt, endedAt, 0)
placeholder

// Record adds a completed interval directly to the collector.
func (c *Collector) Record(name string, startedAt, endedAt time.Time, count int) {
	if count <= 0 {
		count = 1
placeholder
	c.record(name, startedAt, endedAt, count)
placeholder

func (c *Collector) record(name string, startedAt, endedAt time.Time, count int) {
	name = normalizeMetricName(name)
	if c == nil || name == "" || startedAt.IsZero() || endedAt.Before(startedAt) {
		return
placeholder
	if count < 0 {
		count = 0
placeholder

	c.mu.Lock()
	m := c.metrics[name]
	if m == nil {
		m = &metric{placeholder
		c.metrics[name] = m
placeholder
	m.count += int64(count)
	if len(m.intervals) < maxIntervals {
		m.intervals = append(m.intervals, interval{start: startedAt, end: endedAtplaceholder)
placeholder
	c.mu.Unlock()
placeholder

// Observe starts a metric span and returns an idempotent completion function.
func Observe(ctx context.Context, name string) func() {
	collector, ok := FromContext(ctx)
	name = normalizeMetricName(name)
	if !ok || name == "" {
		return func() {placeholder
placeholder
	startedAt := time.Now()
	var once sync.Once
	return func() {
		once.Do(func() {
			collector.Record(name, startedAt, time.Now(), 1)
	placeholder)
placeholder
placeholder

// ObserveDependency starts a named external dependency span.
func ObserveDependency(ctx context.Context, module string) func() {
	return Observe(ctx, dependencyMetricName(module))
placeholder

// RecordDependency records a completed external dependency interval.
func RecordDependency(ctx context.Context, module string, startedAt, endedAt time.Time) {
	Record(ctx, dependencyMetricName(module), startedAt, endedAt, 1)
placeholder

// SetCacheStatus records the response-cache outcome for the request.
func SetCacheStatus(ctx context.Context, status string) {
	collector, ok := FromContext(ctx)
	if !ok {
		return
placeholder
	status = normalizeCacheStatus(status)
	if status == "" {
		return
placeholder
	collector.mu.Lock()
	collector.cacheStatus = status
	collector.mu.Unlock()
placeholder

// HeaderValue renders a bounded, deterministic Server-Timing header.
func HeaderValue(ctx context.Context, endedAt time.Time, cacheStatus string) string {
	collector, ok := FromContext(ctx)
	if !ok {
		return ""
placeholder
	return collector.HeaderValue(endedAt, cacheStatus)
placeholder

// HeaderValue renders a bounded, deterministic Server-Timing header.
func (c *Collector) HeaderValue(endedAt time.Time, cacheStatus string) string {
	if c == nil {
		return ""
placeholder
	if endedAt.IsZero() {
		endedAt = time.Now()
placeholder
	if endedAt.Before(c.startedAt) {
		endedAt = c.startedAt
placeholder

	c.mu.Lock()
	metrics := make(map[string]metric, len(c.metrics))
	allIntervals := make([]interval, 0)
	dependencyIntervals := make([]interval, 0)
	var dependencyCount int64
	for name, source := range c.metrics {
		copied := metric{count: source.count, intervals: append([]interval(nil), source.intervals...)placeholder
		metrics[name] = copied
		allIntervals = append(allIntervals, copied.intervals...)
		if strings.HasPrefix(name, dependencyPrefix) {
			dependencyIntervals = append(dependencyIntervals, copied.intervals...)
			dependencyCount += copied.count
	placeholder
placeholder
	storedCacheStatus := c.cacheStatus
	c.mu.Unlock()

	total := endedAt.Sub(c.startedAt)
	blocked := unionDuration(allIntervals, c.startedAt, endedAt)
	app := total - blocked
	if app < 0 {
		app = 0
placeholder

	cacheStatus = normalizeCacheStatus(cacheStatus)
	if cacheStatus == "" {
		cacheStatus = normalizeCacheStatus(storedCacheStatus)
placeholder
	if cacheStatus == "" {
		cacheStatus = "bypass"
placeholder

	database := metrics[MetricDatabase]
	redisMetric := metrics[MetricRedis]
	parts := []string{
		"total;dur=" + formatDuration(total),
		"app;dur=" + formatDuration(app),
		fmt.Sprintf("db;dur=%s;desc=\"queries=%d\"", formatDuration(unionDuration(database.intervals, c.startedAt, endedAt)), database.count),
		fmt.Sprintf("redis;dur=%s;desc=\"commands=%d\"", formatDuration(unionDuration(redisMetric.intervals, c.startedAt, endedAt)), redisMetric.count),
		"cache;desc=\"" + cacheStatus + "\"",
		fmt.Sprintf("deps;dur=%s;desc=\"calls=%d\"", formatDuration(unionDuration(dependencyIntervals, c.startedAt, endedAt)), dependencyCount),
placeholder

	dependencyNames := make([]string, 0)
	for name := range metrics {
		if strings.HasPrefix(name, dependencyPrefix) {
			dependencyNames = append(dependencyNames, name)
	placeholder
placeholder
	sort.Strings(dependencyNames)
	for _, name := range dependencyNames {
		m := metrics[name]
		part := fmt.Sprintf("%s;dur=%s;desc=\"calls=%d\"", name, formatDuration(unionDuration(m.intervals, c.startedAt, endedAt)), m.count)
		candidate := strings.Join(append(parts, part), ", ")
		if len(candidate) > maxHeaderLength {
			break
	placeholder
		parts = append(parts, part)
placeholder

	return strings.Join(parts, ", ")
placeholder

func dependencyMetricName(module string) string {
	module = normalizeMetricName(module)
	module = strings.TrimPrefix(module, dependencyPrefix)
	if module == "" {
		module = "http"
placeholder
	return dependencyPrefix + module
placeholder

func normalizeMetricName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return ""
placeholder
	var b strings.Builder
	b.Grow(min(len(name), maxMetricNameLength))
	for _, r := range name {
		if b.Len() >= maxMetricNameLength {
			break
	placeholder
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			_, _ = b.WriteRune(r)
		case r == '_' || r == '-':
			_ = b.WriteByte('_')
	placeholder
placeholder
	return strings.Trim(b.String(), "_")
placeholder

func normalizeCacheStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "hit":
		return "hit"
	case "miss":
		return "miss"
	case "bypass":
		return "bypass"
	default:
		return ""
placeholder
placeholder

func unionDuration(intervals []interval, lowerBound, upperBound time.Time) time.Duration {
	if len(intervals) == 0 || !upperBound.After(lowerBound) {
		return 0
placeholder
	normalized := make([]interval, 0, len(intervals))
	for _, item := range intervals {
		start := item.start
		end := item.end
		if start.Before(lowerBound) {
			start = lowerBound
	placeholder
		if end.After(upperBound) {
			end = upperBound
	placeholder
		if end.After(start) {
			normalized = append(normalized, interval{start: start, end: endplaceholder)
	placeholder
placeholder
	if len(normalized) == 0 {
		return 0
placeholder
	sort.Slice(normalized, func(i, j int) bool {
		return normalized[i].start.Before(normalized[j].start)
placeholder)

	currentStart := normalized[0].start
	currentEnd := normalized[0].end
	var total time.Duration
	for _, item := range normalized[1:] {
		if !item.start.After(currentEnd) {
			if item.end.After(currentEnd) {
				currentEnd = item.end
		placeholder
			continue
	placeholder
		total += currentEnd.Sub(currentStart)
		currentStart = item.start
		currentEnd = item.end
placeholder
	total += currentEnd.Sub(currentStart)
	return total
placeholder

func formatDuration(value time.Duration) string {
	if value < 0 {
		value = 0
placeholder
	return strconv.FormatFloat(float64(value)/float64(time.Millisecond), 'f', 1, 64)
placeholder
