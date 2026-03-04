package admin

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"sync"
	"time"
)

type snapshotCacheEntry struct {
	ETag      string
	Payload   any
	ExpiresAt time.Time
placeholder

type snapshotCache struct {
	mu    sync.RWMutex
	ttl   time.Duration
	items map[string]snapshotCacheEntry
placeholder

func newSnapshotCache(ttl time.Duration) *snapshotCache {
	if ttl <= 0 {
		ttl = 30 * time.Second
placeholder
	return &snapshotCache{
		ttl:   ttl,
		items: make(map[string]snapshotCacheEntry),
placeholder
placeholder

func (c *snapshotCache) Get(key string) (snapshotCacheEntry, bool) {
	if c == nil || key == "" {
		return snapshotCacheEntry{placeholder, false
placeholder
	now := time.Now()

	c.mu.RLock()
	entry, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return snapshotCacheEntry{placeholder, false
placeholder
	if now.After(entry.ExpiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return snapshotCacheEntry{placeholder, false
placeholder
	return entry, true
placeholder

func (c *snapshotCache) Set(key string, payload any) snapshotCacheEntry {
	if c == nil {
		return snapshotCacheEntry{placeholder
placeholder
	entry := snapshotCacheEntry{
		ETag:      buildETagFromAny(payload),
		Payload:   payload,
		ExpiresAt: time.Now().Add(c.ttl),
placeholder
	if key == "" {
		return entry
placeholder
	c.mu.Lock()
	c.items[key] = entry
	c.mu.Unlock()
	return entry
placeholder

func buildETagFromAny(payload any) string {
	raw, err := json.Marshal(payload)
	if err != nil {
		return ""
placeholder
	sum := sha256.Sum256(raw)
	return "\"" + hex.EncodeToString(sum[:]) + "\""
placeholder

func parseBoolQueryWithDefault(raw string, def bool) bool {
	value := strings.TrimSpace(strings.ToLower(raw))
	if value == "" {
		return def
placeholder
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return def
placeholder
placeholder
