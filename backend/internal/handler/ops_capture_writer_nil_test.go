package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type deterministicOpsCaptureWriterStatePool struct {
	states []*opsCaptureWriterState
placeholder

func (p *deterministicOpsCaptureWriterStatePool) Get() any {
	if len(p.states) == 0 {
		return &opsCaptureWriterState{limit: opsCaptureWriterLimitplaceholder
placeholder
	last := len(p.states) - 1
	state := p.states[last]
	p.states = p.states[:last]
	return state
placeholder

func (p *deterministicOpsCaptureWriterStatePool) Put(value any) {
	if state, ok := value.(*opsCaptureWriterState); ok && state != nil {
		p.states = append(p.states, state)
placeholder
placeholder

func TestOpsCaptureWriter_NilInnerWriter_NoPanic(t *testing.T) {
	w := &opsCaptureWriter{placeholder

	assert.NotPanics(t, func() {
		assert.Equal(t, 0, w.Status())
placeholder)
	assert.NotPanics(t, func() {
		assert.Equal(t, -1, w.Size())
placeholder)
	assert.NotPanics(t, func() {
		assert.False(t, w.Written())
placeholder)
	assert.NotPanics(t, func() {
		n, err := w.Write([]byte("test"))
		assert.Equal(t, 0, n)
		assert.NoError(t, err)
placeholder)
	assert.NotPanics(t, func() {
		n, err := w.WriteString("test")
		assert.Equal(t, 0, n)
		assert.NoError(t, err)
placeholder)
	assert.NotPanics(t, func() {
		h := w.Header()
		assert.NotNil(t, h)
placeholder)
	assert.NotPanics(t, func() {
		w.WriteHeader(200)
placeholder)
	assert.NotPanics(t, func() {
		w.WriteHeaderNow()
placeholder)
	assert.NotPanics(t, func() {
		w.Flush()
placeholder)
	assert.NotPanics(t, func() {
		conn, rw, err := w.Hijack()
		assert.Nil(t, conn)
		assert.Nil(t, rw)
		assert.Error(t, err)
placeholder)
	assert.NotPanics(t, func() {
		ch := w.CloseNotify()
		assert.NotNil(t, ch)
placeholder)
	assert.NotPanics(t, func() {
		p := w.Pusher()
		assert.Nil(t, p)
placeholder)
placeholder

func TestOpsCaptureWriter_CompactKeepaliveRestoresOriginalWriter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	outerStatus := -1
	router.Use(func(c *gin.Context) {
		c.Next()
		outerStatus = c.Writer.Status()
placeholder)
	router.Use(OpsErrorLoggerMiddleware(nil))
	router.GET("/compact", func(c *gin.Context) {
		service.MarkOpenAICompactClientStream(c)
		stop := service.StartOpenAICompactSSEKeepalive(c, time.Hour)
		defer stop()
		c.Status(http.StatusOK)
placeholder)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/compact", nil)
	require.NotPanics(t, func() {
		router.ServeHTTP(recorder, request)
placeholder)
	require.Equal(t, http.StatusOK, outerStatus)
	require.Equal(t, http.StatusOK, recorder.Code)
placeholder

func TestOpsCaptureWriter_StaleLeaseCannotReachReacquiredState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := &deterministicOpsCaptureWriterStatePool{placeholder

	firstRecorder := httptest.NewRecorder()
	firstContext, _ := gin.CreateTestContext(firstRecorder)
	stale := acquireOpsCaptureWriterFromPool(pool, firstContext.Writer)
	releaseOpsCaptureWriter(stale)

	secondRecorder := httptest.NewRecorder()
	secondContext, _ := gin.CreateTestContext(secondRecorder)
	current := acquireOpsCaptureWriterFromPool(pool, secondContext.Writer)
	defer releaseOpsCaptureWriter(current)
	require.NotSame(t, stale, current)
	require.Same(t, stale.state, current.state)

	current.WriteHeader(http.StatusInternalServerError)
	_, err := current.WriteString("current")
placeholder
	require.Equal(t, []byte("current"), current.capturedBytes())

	n, err := stale.WriteString("stale")
placeholder
	require.Zero(t, n)
	require.Nil(t, stale.capturedBytes())
	require.Equal(t, []byte("current"), current.capturedBytes())
	require.NotContains(t, secondRecorder.Body.String(), "stale")

	// Releasing the stale handle must not return an active state to the pool.
	releaseOpsCaptureWriter(stale)
	thirdRecorder := httptest.NewRecorder()
	thirdContext, _ := gin.CreateTestContext(thirdRecorder)
	other := acquireOpsCaptureWriterFromPool(pool, thirdContext.Writer)
	defer releaseOpsCaptureWriter(other)
	require.NotSame(t, current.state, other.state)
placeholder
