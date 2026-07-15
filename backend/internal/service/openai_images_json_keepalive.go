package service

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const openAIImagesJSONKeepaliveKey = "openai_images_json_keepalive"

// openAIImagesJSONKeepalive keeps non-streaming Images API requests alive while
// an OAuth upstream is producing SSE internally. JSON permits leading
// whitespace, so each heartbeat remains compatible with clients expecting one
// final JSON document.
//
// Once the first heartbeat is sent, the HTTP status is committed as 200. Late
// upstream errors are still returned as an OpenAI-compatible JSON error body,
// matching the status tradeoff used by the compact SSE keepalive path.
type openAIImagesJSONKeepalive struct {
	mu      sync.Mutex
	writer  gin.ResponseWriter
	started bool
	stopped bool
	bytes   int
	stop    chan struct{placeholder
placeholder

// StartOpenAIImagesJSONKeepalive starts whitespace heartbeats for a
// non-streaming Images request. A non-positive interval disables the feature.
func StartOpenAIImagesJSONKeepalive(c *gin.Context, interval time.Duration) func() {
	if c == nil || c.Writer == nil || interval <= 0 {
		return func() {placeholder
placeholder
	originalWriter := c.Writer
	k := &openAIImagesJSONKeepalive{
		writer: originalWriter,
		stop:   make(chan struct{placeholder),
placeholder
	c.Set(openAIImagesJSONKeepaliveKey, k)
	wrappedWriter := &openAIImagesJSONKeepaliveWriter{ResponseWriter: originalWriter, k: kplaceholder
	c.Writer = wrappedWriter

	var reqDone <-chan struct{placeholder
	if c.Request != nil {
		reqDone = c.Request.Context().Done()
placeholder
	go func() {
		timer := time.NewTimer(interval)
		defer timer.Stop()
		for {
			select {
			case <-k.stop:
				return
			case <-reqDone:
				return
			case <-timer.C:
		placeholder
			if !k.beat() {
				return
		placeholder
			timer.Reset(interval)
	placeholder
placeholder()

	return func() {
		k.Stop()
		if current, ok := c.Writer.(*openAIImagesJSONKeepaliveWriter); ok && current == wrappedWriter {
			c.Writer = originalWriter
	placeholder
placeholder
placeholder

func (k *openAIImagesJSONKeepalive) beat() bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.stopped {
		return false
placeholder
	if !k.started {
		header := k.writer.Header()
		header.Set("Content-Type", "application/json; charset=utf-8")
		header.Set("Cache-Control", "no-cache")
		header.Set("X-Accel-Buffering", "no")
		k.writer.WriteHeader(http.StatusOK)
		k.started = true
placeholder
	n, err := k.writer.Write([]byte(" \n"))
	k.bytes += n
	if err != nil {
		k.stopped = true
		return false
placeholder
	k.writer.Flush()
	return true
placeholder

func (k *openAIImagesJSONKeepalive) Stop() {
	k.mu.Lock()
	k.markStoppedLocked()
	k.mu.Unlock()
placeholder

func (k *openAIImagesJSONKeepalive) markStoppedLocked() {
	if k.stopped {
		return
placeholder
	k.stopped = true
	close(k.stop)
placeholder

// StopOpenAIImagesJSONKeepaliveCommitted stops heartbeats and reports whether
// they already committed a 200 response.
func StopOpenAIImagesJSONKeepaliveCommitted(c *gin.Context) bool {
	k := openAIImagesJSONKeepaliveFromContext(c)
	if k == nil {
		return false
placeholder
	k.mu.Lock()
	k.markStoppedLocked()
	committed := k.started
	k.mu.Unlock()
	return committed
placeholder

// OpenAIImagesJSONKeepalivePresent reports whether the response writer belongs
// to an Images JSON request, including fast responses before the first beat.
func OpenAIImagesJSONKeepalivePresent(c *gin.Context) bool {
	return openAIImagesJSONKeepaliveFromContext(c) != nil
placeholder

// OpenAIImagesJSONKeepaliveAdjustedWrittenSize excludes heartbeat whitespace
// from response-size checks so account retry and failover remain available.
func OpenAIImagesJSONKeepaliveAdjustedWrittenSize(c *gin.Context) int {
	if c == nil || c.Writer == nil {
		return -1
placeholder
	k := openAIImagesJSONKeepaliveFromContext(c)
	if k == nil {
		return c.Writer.Size()
placeholder
	k.mu.Lock()
	defer k.mu.Unlock()
	size := k.writer.Size()
	if size < 0 {
		return size
placeholder
	if real := size - k.bytes; real > 0 {
		return real
placeholder
	return -1
placeholder

func openAIImagesJSONKeepaliveFromContext(c *gin.Context) *openAIImagesJSONKeepalive {
	if c == nil {
		return nil
placeholder
	value, ok := c.Get(openAIImagesJSONKeepaliveKey)
	if !ok {
		return nil
placeholder
	k, _ := value.(*openAIImagesJSONKeepalive)
	return k
placeholder

type openAIImagesJSONKeepaliveWriter struct {
	gin.ResponseWriter
	k *openAIImagesJSONKeepalive
placeholder

func (w *openAIImagesJSONKeepaliveWriter) suspend() {
	if w.k != nil {
		w.k.Stop()
placeholder
placeholder

func (w *openAIImagesJSONKeepaliveWriter) Header() http.Header {
	w.suspend()
	if w.ResponseWriter == nil {
		return http.Header{placeholder
placeholder
	return w.ResponseWriter.Header()
placeholder

func (w *openAIImagesJSONKeepaliveWriter) Write(data []byte) (int, error) {
	w.suspend()
	if w.ResponseWriter == nil {
		return 0, nil
placeholder
	return w.ResponseWriter.Write(data)
placeholder

func (w *openAIImagesJSONKeepaliveWriter) WriteString(s string) (int, error) {
	w.suspend()
	if w.ResponseWriter == nil {
		return 0, nil
placeholder
	return w.ResponseWriter.WriteString(s)
placeholder

func (w *openAIImagesJSONKeepaliveWriter) WriteHeader(code int) {
	w.suspend()
	if w.ResponseWriter != nil {
		w.ResponseWriter.WriteHeader(code)
placeholder
placeholder

func (w *openAIImagesJSONKeepaliveWriter) WriteHeaderNow() {
	w.suspend()
	if w.ResponseWriter != nil {
		w.ResponseWriter.WriteHeaderNow()
placeholder
placeholder

func (w *openAIImagesJSONKeepaliveWriter) Flush() {
	w.suspend()
	if w.ResponseWriter != nil {
		w.ResponseWriter.Flush()
placeholder
placeholder

func (w *openAIImagesJSONKeepaliveWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.ResponseWriter == nil {
		return nil, nil, errors.New("response writer released")
placeholder
	return w.ResponseWriter.Hijack()
placeholder

func (w *openAIImagesJSONKeepaliveWriter) CloseNotify() <-chan bool {
	if w.ResponseWriter == nil {
		ch := make(chan bool)
		close(ch)
		return ch
placeholder
	return w.ResponseWriter.CloseNotify()
placeholder

func (w *openAIImagesJSONKeepaliveWriter) Pusher() http.Pusher {
	if w.ResponseWriter == nil {
		return nil
placeholder
	return w.ResponseWriter.Pusher()
placeholder

func (w *openAIImagesJSONKeepaliveWriter) Status() int {
	if w.k == nil || w.ResponseWriter == nil {
		return 0
placeholder
	w.k.mu.Lock()
	defer w.k.mu.Unlock()
	return w.ResponseWriter.Status()
placeholder

func (w *openAIImagesJSONKeepaliveWriter) Size() int {
	if w.k == nil || w.ResponseWriter == nil {
		return 0
placeholder
	w.k.mu.Lock()
	defer w.k.mu.Unlock()
	return w.ResponseWriter.Size()
placeholder

func (w *openAIImagesJSONKeepaliveWriter) Written() bool {
	if w.k == nil || w.ResponseWriter == nil {
		return false
placeholder
	w.k.mu.Lock()
	defer w.k.mu.Unlock()
	return w.ResponseWriter.Written()
placeholder
