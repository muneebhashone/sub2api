package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var errOpenAIWSSessionPreempted = errors.New("openai ws session preempted by newer request")

const (
	openAIWSSessionPreemptOwnerTTL      = 2 * time.Hour
	openAIWSSessionPreemptWatchInterval = 2 * time.Second
	openAIWSSessionPreemptCachePrefix   = "wspreempt:"
)

// OpenAIWSSessionPreemptionCache is an optional GatewayCache capability. The
// production Redis cache implements all operations atomically; cache stubs do
// not need to implement it for ordinary gateway tests.
type OpenAIWSSessionPreemptionCache interface {
	ClaimOpenAIResponsesSessionWindow(ctx context.Context, groupID int64, sessionHash string, owner []byte, ttl time.Duration) ([]byte, error)
	CompareAndRefreshOpenAIResponsesSessionWindow(ctx context.Context, groupID int64, sessionHash string, expected []byte, ttl time.Duration) (bool, error)
	CompareAndDeleteOpenAIResponsesSessionWindow(ctx context.Context, groupID int64, sessionHash string, expected []byte) (bool, error)
placeholder

func NewOpenAIWSSessionPreemptedError() error {
	return errOpenAIWSSessionPreempted
placeholder

type openAIWSSessionPreemptKey struct {
	groupID     int64
	apiKeyID    int64
	sessionHash string
placeholder

func newOpenAIWSSessionPreemptKey(groupID, apiKeyID int64, sessionHash string) (openAIWSSessionPreemptKey, bool) {
	sessionHash = strings.TrimSpace(sessionHash)
	if groupID <= 0 || apiKeyID <= 0 || sessionHash == "" {
		return openAIWSSessionPreemptKey{placeholder, false
placeholder
	return openAIWSSessionPreemptKey{groupID: groupID, apiKeyID: apiKeyID, sessionHash: sessionHashplaceholder, true
placeholder

func openAIWSSessionPreemptCacheHash(apiKeyID int64, sessionHash string) string {
	return fmt.Sprintf("%s%d:%s", openAIWSSessionPreemptCachePrefix, apiKeyID, strings.TrimSpace(sessionHash))
placeholder

type openAIWSSessionPreemptEntry struct {
	generation uint64
	cancel     func()
placeholder

type openAIWSSessionPreemptRegistry struct {
	mu     sync.Mutex
	next   uint64
	active map[openAIWSSessionPreemptKey]openAIWSSessionPreemptEntry
placeholder

func (r *openAIWSSessionPreemptRegistry) Begin(key openAIWSSessionPreemptKey, cancel func()) (cleanup func(), preemptedPrevious bool) {
	if r == nil || strings.TrimSpace(key.sessionHash) == "" {
		return func() {placeholder, false
placeholder
	r.mu.Lock()
	if r.active == nil {
		r.active = make(map[openAIWSSessionPreemptKey]openAIWSSessionPreemptEntry)
placeholder
	r.next++
	generation := r.next
	previous, hadPrevious := r.active[key]
	r.active[key] = openAIWSSessionPreemptEntry{generation: generation, cancel: cancelplaceholder
	r.mu.Unlock()
	if hadPrevious && previous.cancel != nil {
		previous.cancel()
placeholder
	return func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		current, ok := r.active[key]
		if ok && current.generation == generation {
			delete(r.active, key)
	placeholder
placeholder, hadPrevious
placeholder

func (s *OpenAIGatewayService) beginOpenAIWSSessionPreemptContext(
	ctx context.Context,
	account *Account,
	groupID, apiKeyID int64,
	sessionHash string,
	httpIngressWSOneShot bool,
) (context.Context, func(), bool, bool) {
	if ctx == nil {
		ctx = context.Background()
placeholder
	if s == nil || account == nil || account.Platform != PlatformOpenAI || account.Type != AccountTypeOAuth || httpIngressWSOneShot {
		return ctx, func() {placeholder, false, false
placeholder
	key, ok := newOpenAIWSSessionPreemptKey(groupID, apiKeyID, sessionHash)
	if !ok {
		return ctx, func() {placeholder, false, false
placeholder

	preemptCtx, cancel := context.WithCancelCause(ctx)
	ownerToken := uuid.NewString()
	var preemptOnce sync.Once
	preempt := func() {
		preemptOnce.Do(func() {
			if stateStore := s.getOpenAIWSStateStore(); stateStore != nil {
				stateStore.DeleteSessionTurnState(key.groupID, key.sessionHash)
				stateStore.DeleteSessionConn(key.groupID, key.sessionHash)
		placeholder
			cancel(errOpenAIWSSessionPreempted)
	placeholder)
placeholder
	previousRemoteOwner, remoteClaimed := s.claimOpenAIWSSessionPreemptOwner(ctx, key, ownerToken)
	preemptedPrevious := remoteClaimed && previousRemoteOwner != "" && previousRemoteOwner != ownerToken
	cleanupLocal, hadLocalPrevious := s.openaiWSSessionPreemptions.Begin(key, preempt)
	preemptedPrevious = preemptedPrevious || hadLocalPrevious
	stopWatch := func() {placeholder
	if remoteClaimed {
		stopWatch = s.watchOpenAIWSSessionPreemptOwner(preemptCtx, key, ownerToken, preempt)
placeholder

	return preemptCtx, func() {
		stopWatch()
		cleanupLocal()
		if remoteClaimed {
			s.releaseOpenAIWSSessionPreemptOwner(context.Background(), key, ownerToken)
	placeholder
		cancel(nil)
placeholder, true, preemptedPrevious
placeholder

func (s *OpenAIGatewayService) openAIWSSessionPreemptionCache() OpenAIWSSessionPreemptionCache {
	if s == nil || s.cache == nil {
		return nil
placeholder
	cache, _ := s.cache.(OpenAIWSSessionPreemptionCache)
	return cache
placeholder

func (s *OpenAIGatewayService) claimOpenAIWSSessionPreemptOwner(ctx context.Context, key openAIWSSessionPreemptKey, ownerToken string) (string, bool) {
	cache := s.openAIWSSessionPreemptionCache()
	if cache == nil || strings.TrimSpace(ownerToken) == "" {
		return "", false
placeholder
	cacheCtx, cancel := context.WithTimeout(ctx, openAIWSStateStoreRedisTimeout)
	defer cancel()
	previous, err := cache.ClaimOpenAIResponsesSessionWindow(
		cacheCtx,
		key.groupID,
		openAIWSSessionPreemptCacheHash(key.apiKeyID, key.sessionHash),
		[]byte(strings.TrimSpace(ownerToken)),
		openAIWSSessionPreemptOwnerTTL,
	)
	if err != nil {
		return "", false
placeholder
	return strings.TrimSpace(string(previous)), true
placeholder

func (s *OpenAIGatewayService) releaseOpenAIWSSessionPreemptOwner(ctx context.Context, key openAIWSSessionPreemptKey, ownerToken string) {
	cache := s.openAIWSSessionPreemptionCache()
	if cache == nil || strings.TrimSpace(ownerToken) == "" {
		return
placeholder
	cacheCtx, cancel := context.WithTimeout(ctx, openAIWSStateStoreRedisTimeout)
	defer cancel()
	_, _ = cache.CompareAndDeleteOpenAIResponsesSessionWindow(
		cacheCtx,
		key.groupID,
		openAIWSSessionPreemptCacheHash(key.apiKeyID, key.sessionHash),
		[]byte(strings.TrimSpace(ownerToken)),
	)
placeholder

func (s *OpenAIGatewayService) watchOpenAIWSSessionPreemptOwner(ctx context.Context, key openAIWSSessionPreemptKey, ownerToken string, onLost func()) func() {
	cache := s.openAIWSSessionPreemptionCache()
	if cache == nil || onLost == nil || strings.TrimSpace(ownerToken) == "" {
		return func() {placeholder
placeholder
	stopCh := make(chan struct{placeholder)
	var once sync.Once
	go func() {
		ticker := time.NewTicker(openAIWSSessionPreemptWatchInterval)
		defer ticker.Stop()
		for {
			select {
			case <-stopCh:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				cacheCtx, cancel := context.WithTimeout(context.Background(), openAIWSStateStoreRedisTimeout)
				owned, err := cache.CompareAndRefreshOpenAIResponsesSessionWindow(
					cacheCtx,
					key.groupID,
					openAIWSSessionPreemptCacheHash(key.apiKeyID, key.sessionHash),
					[]byte(strings.TrimSpace(ownerToken)),
					openAIWSSessionPreemptOwnerTTL,
				)
				cancel()
				if err == nil && !owned {
					onLost()
					return
			placeholder
		placeholder
	placeholder
placeholder()
	return func() { once.Do(func() { close(stopCh) placeholder) placeholder
placeholder

func isOpenAIWSSessionPreempted(ctx context.Context) bool {
	return ctx != nil && errors.Is(context.Cause(ctx), errOpenAIWSSessionPreempted)
placeholder

func IsOpenAIWSSessionPreemptedError(err error) bool {
	if err == nil {
		return false
placeholder
	if errors.Is(err, errOpenAIWSSessionPreempted) {
		return true
placeholder
	var fallbackErr *openAIWSFallbackError
	return errors.As(err, &fallbackErr) && fallbackErr != nil && strings.TrimPrefix(strings.TrimSpace(fallbackErr.Reason), "prewarm_") == "session_preempted"
placeholder
