package handler

import (
	"context"
	"sync"
	"time"
)

type imageConcurrencyLimiter struct {
	mu      sync.Mutex
	notify  chan struct{placeholder
	limit   int
	active  int
	waiting int
	enabled bool
placeholder

func (l *imageConcurrencyLimiter) TryAcquire(enabled bool, limit int) (func(), bool) {
	return l.acquire(context.Background(), enabled, limit, false, 0, 0)
placeholder

func (l *imageConcurrencyLimiter) Acquire(ctx context.Context, enabled bool, limit int, wait bool, timeout time.Duration, maxWaiting int) (func(), bool) {
	return l.acquire(ctx, enabled, limit, wait, timeout, maxWaiting)
placeholder

func (l *imageConcurrencyLimiter) acquire(ctx context.Context, enabled bool, limit int, wait bool, timeout time.Duration, maxWaiting int) (func(), bool) {
	if !enabled || limit <= 0 {
		return nil, true
placeholder
	if ctx == nil {
		ctx = context.Background()
placeholder
	if wait {
		if timeout <= 0 {
			return nil, false
	placeholder
		waitCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		ctx = waitCtx
placeholder
	if maxWaiting < 0 {
		maxWaiting = 0
placeholder
	for {
		release, acquired, waitRelease, notify := l.tryAcquireLocked(enabled, limit, wait, maxWaiting)
		if acquired {
			return release, acquired
	placeholder
		if !wait || notify == nil {
			return nil, false
	placeholder
		if !l.waitForSlot(ctx, notify) {
			if waitRelease != nil {
				waitRelease()
		placeholder
			return nil, false
	placeholder
		if waitRelease != nil {
			waitRelease()
	placeholder
placeholder
placeholder

func (l *imageConcurrencyLimiter) tryAcquireLocked(enabled bool, limit int, wait bool, maxWaiting int) (func(), bool, func(), <-chan struct{placeholder) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.notify == nil {
		l.notify = make(chan struct{placeholder)
placeholder
	if l.enabled != enabled || l.limit != limit {
		l.enabled = enabled
		l.limit = limit
placeholder
	if l.active < l.limit {
		l.active++
		return l.releaseFunc(), true, nil, nil
placeholder
	if !wait {
		return nil, false, nil, nil
placeholder
	if maxWaiting > 0 && l.waiting >= maxWaiting {
		return nil, false, nil, nil
placeholder
	l.waiting++
	return nil, false, l.waiterReleaseFunc(), l.notify
placeholder

func (l *imageConcurrencyLimiter) waitForSlot(ctx context.Context, notify <-chan struct{placeholder) bool {
	select {
	case <-notify:
		return true
	case <-ctx.Done():
		return false
placeholder
placeholder

func (l *imageConcurrencyLimiter) releaseFunc() func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			l.mu.Lock()
			if l.active > 0 {
				l.active--
		placeholder
			if l.notify != nil {
				close(l.notify)
				l.notify = make(chan struct{placeholder)
		placeholder
			l.mu.Unlock()
	placeholder)
placeholder
placeholder

func (l *imageConcurrencyLimiter) waiterReleaseFunc() func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			l.mu.Lock()
			if l.waiting > 0 {
				l.waiting--
		placeholder
			l.mu.Unlock()
	placeholder)
placeholder
placeholder
