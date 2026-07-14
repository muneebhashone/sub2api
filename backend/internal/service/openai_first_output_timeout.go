package service

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	openAIFirstOutputStageMemoryLimit        = 64 * 1024
	openAIFirstOutputStageMaxBytes           = 8 * 1024 * 1024
	openAIFirstOutputScannerFramingAllowance = 64
	openAIFirstOutputGuardQueueSize          = 1
	openAIDefaultStreamQueueSize             = 16
)

var (
	errOpenAIFirstOutputStageLimit   = errors.New("openai first-output staging limit exceeded")
	errOpenAIFirstOutputScannerLimit = errors.New("openai pre-output scanner token limit exceeded")
)

type openAIFirstOutputStage struct {
	limit      int64
	size       int64
	memory     bytes.Buffer
	tempFile   *os.File
	tempPath   string
	createTemp func() (*os.File, error)
	removeFile func(string) error
	memoryOnly bool
	cleanupErr error
	closed     bool
placeholder

func newOpenAIFirstOutputStage(limit int64) *openAIFirstOutputStage {
	if limit < 1 {
		limit = 1
placeholder
	return &openAIFirstOutputStage{
		limit:      limit,
		createTemp: func() (*os.File, error) { return os.CreateTemp("", "sub2api-openai-first-output-*") placeholder,
		removeFile: os.Remove,
		memoryOnly: runtime.GOOS == "windows",
placeholder
placeholder

func newDefaultOpenAIFirstOutputStage() *openAIFirstOutputStage {
	return newOpenAIFirstOutputStage(openAIFirstOutputStageMaxBytes)
placeholder

func openAIFirstOutputEventQueueSize(guardFirstOutput bool) int {
	if guardFirstOutput {
		return openAIFirstOutputGuardQueueSize
placeholder
	return openAIDefaultStreamQueueSize
placeholder

func openAIFirstOutputDynamicScanLines(guardActive *atomic.Bool) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		if err != nil || guardActive == nil || !guardActive.Load() {
			return advance, token, err
	placeholder
		limit := openAIFirstOutputStageMaxBytes + openAIFirstOutputScannerFramingAllowance
		if token != nil {
			if len(token) > limit {
				return 0, nil, errOpenAIFirstOutputScannerLimit
		placeholder
			return advance, token, nil
	placeholder
		// At the limit with no delimiter, another byte would necessarily exceed
		// the guarded token budget. Fail before Scanner grows toward MaxLineSize.
		if len(data) >= limit {
			return 0, nil, errOpenAIFirstOutputScannerLimit
	placeholder
		return advance, token, nil
placeholder
placeholder

func (s *openAIFirstOutputStage) Buffered() int64 {
	if s == nil {
		return 0
placeholder
	return s.size
placeholder

func (s *openAIFirstOutputStage) WriteString(value string) (int, error) {
	if err := s.prepareWrite(len(value)); err != nil {
		return 0, err
placeholder
	var n int
	var err error
	if s.tempFile == nil {
		n, err = s.memory.WriteString(value)
placeholder else {
		n, err = io.WriteString(s.tempFile, value)
placeholder
	s.size += int64(n)
	if err != nil {
		return n, fmt.Errorf("write first-output stage: %w", err)
placeholder
	return n, nil
placeholder

func (s *openAIFirstOutputStage) Write(p []byte) (int, error) {
	if err := s.prepareWrite(len(p)); err != nil {
		return 0, err
placeholder
	var n int
	var err error
	if s.tempFile == nil {
		n, err = s.memory.Write(p)
placeholder else {
		n, err = s.tempFile.Write(p)
placeholder
	s.size += int64(n)
	if err != nil {
		return n, fmt.Errorf("write first-output stage: %w", err)
placeholder
	return n, nil
placeholder

func (s *openAIFirstOutputStage) prepareWrite(incoming int) error {
	if s == nil || s.closed {
		return os.ErrClosed
placeholder
	if int64(incoming) > s.limit-s.size {
		return fmt.Errorf("%w: buffered=%d incoming=%d limit=%d", errOpenAIFirstOutputStageLimit, s.size, incoming, s.limit)
placeholder
	if s.tempFile != nil || s.memoryOnly || s.size+int64(incoming) <= openAIFirstOutputStageMemoryLimit {
		return nil
placeholder
	file, err := s.createTemp()
	if err != nil {
		return fmt.Errorf("create first-output spool: %w", err)
placeholder
	path := file.Name()
	// Unlink before writing any request data. Unix keeps the file descriptor
	// readable, while crashes and SIGKILL cannot leave a named plaintext spool.
	if unlinkErr := s.removeFile(path); unlinkErr != nil {
		closeErr := file.Close()
		removeErr := s.removeFile(path)
		if errors.Is(removeErr, os.ErrNotExist) {
			removeErr = nil
	placeholder
		s.memoryOnly = true
		if removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			s.tempPath = path
	placeholder
		s.cleanupErr = errors.Join(
			s.cleanupErr,
			fmt.Errorf("unlink first-output spool before use: %w", unlinkErr),
			closeErr,
			removeErr,
		)
		return nil
placeholder
	if _, err := file.Write(s.memory.Bytes()); err != nil {
		_ = file.Close()
		return fmt.Errorf("initialize first-output spool: %w", err)
placeholder
	s.tempFile = file
	s.tempPath = path
	s.memory.Reset()
	return nil
placeholder

func (s *openAIFirstOutputStage) CommitTo(dst io.Writer) error {
	if s == nil || s.closed {
		return os.ErrClosed
placeholder
	if s.tempFile == nil {
		if _, err := io.Copy(dst, bytes.NewReader(s.memory.Bytes())); err != nil {
			return err
	placeholder
placeholder else {
		if _, err := s.tempFile.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("seek first-output spool: %w", err)
	placeholder
		if _, err := io.CopyN(dst, s.tempFile, s.size); err != nil {
			return err
	placeholder
placeholder
	if err := s.Close(); err != nil {
		// Delivery succeeded. Preserve cleanup failures for the handler's deferred
		// cleanup/logging pass instead of turning committed bytes into a stream error.
		s.cleanupErr = errors.Join(s.cleanupErr, err)
placeholder
	return nil
placeholder

func (s *openAIFirstOutputStage) Close() error {
	if s == nil {
		return nil
placeholder
	if s.closed && s.tempFile == nil && s.tempPath == "" && s.cleanupErr == nil {
		return nil
placeholder
	s.closed = true
	s.size = 0
	s.memory.Reset()
	closeErr := s.cleanupErr
	s.cleanupErr = nil
	if s.tempFile != nil {
		closeErr = errors.Join(closeErr, s.tempFile.Close())
		s.tempFile = nil
placeholder
	if s.tempPath != "" {
		removeErr := s.removeFile(s.tempPath)
		if removeErr == nil || errors.Is(removeErr, os.ErrNotExist) {
			s.tempPath = ""
	placeholder else {
			closeErr = errors.Join(closeErr, removeErr)
	placeholder
placeholder
	return closeErr
placeholder

func (s *OpenAIGatewayService) openAIFirstOutputTimeout(reasoningEffort string) time.Duration {
	if s == nil || s.cfg == nil || s.cfg.Gateway.OpenAIFirstOutputTimeoutSeconds <= 0 {
		return 0
placeholder
	seconds := s.cfg.Gateway.OpenAIFirstOutputTimeoutSeconds
	switch strings.ToLower(strings.TrimSpace(reasoningEffort)) {
	case "high", "xhigh", "max":
		if override := s.cfg.Gateway.OpenAIHighEffortFirstOutputTimeoutSeconds; override > 0 {
			seconds = override
	placeholder
placeholder
	return time.Duration(seconds) * time.Second
placeholder

func (s *OpenAIGatewayService) newOpenAIFirstOutputTimeoutError(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	startTime time.Time,
	originalModel string,
	reasoningEffort string,
	timeout time.Duration,
	phase string,
	responseHeaders http.Header,
) *UpstreamFailoverError {
	elapsed := time.Since(startTime)
	logger.LegacyPrintf(
		"service.openai_gateway",
		"OpenAI first output timeout: account=%d model=%s effort=%s phase=%s elapsed=%s limit=%s",
		account.ID, originalModel, reasoningEffort, phase, elapsed, timeout,
	)
	requestID := strings.TrimSpace(responseHeaders.Get("x-request-id"))
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Platform: account.Platform, AccountID: account.ID, AccountName: account.Name,
		UpstreamStatusCode: http.StatusGatewayTimeout, UpstreamRequestID: requestID,
		Kind: "first_output_timeout", Message: "OpenAI upstream produced no semantic output before the deadline",
		Detail: fmt.Sprintf("phase=%s elapsed_ms=%d timeout_ms=%d", phase, elapsed.Milliseconds(), timeout.Milliseconds()),
placeholder)
	if s.rateLimitService != nil {
		s.rateLimitService.HandleStreamTimeout(ctx, account, originalModel)
placeholder
	return &UpstreamFailoverError{
		StatusCode:      http.StatusGatewayTimeout,
		ResponseBody:    []byte(`{"error":{"type":"first_output_timeout","message":"Upstream produced no output before the deadline"placeholderplaceholder`),
		ResponseHeaders: responseHeaders.Clone(), SafeToFailoverAfterWrite: true,
placeholder
placeholder

type openAIFirstOutputHeaderGuard struct {
	cancel  context.CancelFunc
	release context.CancelFunc
	timer   *time.Timer
	fired   chan struct{placeholder
	once    sync.Once
placeholder

func newOpenAIFirstOutputHeaderGuard(
	ctx context.Context,
	release context.CancelFunc,
	deadline time.Time,
) (context.Context, *openAIFirstOutputHeaderGuard) {
	guardedCtx, cancel := context.WithCancel(ctx)
	guard := &openAIFirstOutputHeaderGuard{cancel: cancel, release: release, fired: make(chan struct{placeholder)placeholder
	remaining := time.Until(deadline)
	if remaining <= 0 {
		remaining = time.Nanosecond
placeholder
	guard.timer = time.AfterFunc(remaining, func() {
		close(guard.fired)
		cancel()
placeholder)
	return guardedCtx, guard
placeholder

func (g *openAIFirstOutputHeaderGuard) stopHeaderWait() bool {
	if g.timer.Stop() {
		return false
placeholder
	<-g.fired
	return true
placeholder

func (g *openAIFirstOutputHeaderGuard) close() {
	g.once.Do(func() {
		g.timer.Stop()
		g.cancel()
		g.release()
placeholder)
placeholder

type openAIRequestContextReadCloser struct {
	io.ReadCloser
	cleanup func()
	once    sync.Once
	err     error
placeholder

func (r *openAIRequestContextReadCloser) Close() error {
	r.once.Do(func() {
		r.cleanup()
		r.err = r.ReadCloser.Close()
placeholder)
	return r.err
placeholder
