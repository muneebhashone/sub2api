package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpsCaptureWriter_NilInnerWriter_NoPanic(t *testing.T) {
	w := &opsCaptureWriter{placeholder
	w.ResponseWriter = nil

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
