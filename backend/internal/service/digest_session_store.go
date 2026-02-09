package service

import (
	"strconv"
	"strings"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// digestSessionTTL 摘要会话默认 TTL
const digestSessionTTL = 5 * time.Minute

// sessionEntry flat cache 条目
type sessionEntry struct {
	uuid      string
	accountID int64
placeholder

// DigestSessionStore 内存摘要会话存储（flat cache 实现）
// key: "{groupIDplaceholder:{prefixHashplaceholder|{digestChainplaceholder" → *sessionEntry
type DigestSessionStore struct {
	cache *gocache.Cache
placeholder

// NewDigestSessionStore 创建内存摘要会话存储
func NewDigestSessionStore() *DigestSessionStore {
	return &DigestSessionStore{
		cache: gocache.New(digestSessionTTL, time.Minute),
placeholder
placeholder

// Save 保存摘要会话。oldDigestChain 为 Find 返回的 matchedChain，用于删旧 key。
func (s *DigestSessionStore) Save(groupID int64, prefixHash, digestChain, uuid string, accountID int64, oldDigestChain string) {
	if digestChain == "" {
		return
placeholder
	ns := buildNS(groupID, prefixHash)
	s.cache.Set(ns+digestChain, &sessionEntry{uuid: uuid, accountID: accountIDplaceholder, gocache.DefaultExpiration)
	if oldDigestChain != "" && oldDigestChain != digestChain {
		s.cache.Delete(ns + oldDigestChain)
placeholder
placeholder

// Find 查找摘要会话，从完整 chain 逐段截断，返回最长匹配及对应 matchedChain。
func (s *DigestSessionStore) Find(groupID int64, prefixHash, digestChain string) (uuid string, accountID int64, matchedChain string, found bool) {
	if digestChain == "" {
		return "", 0, "", false
placeholder
	ns := buildNS(groupID, prefixHash)
	chain := digestChain
	for {
		if val, ok := s.cache.Get(ns + chain); ok {
			if e, ok := val.(*sessionEntry); ok {
				return e.uuid, e.accountID, chain, true
		placeholder
	placeholder
		i := strings.LastIndex(chain, "-")
		if i < 0 {
			return "", 0, "", false
	placeholder
		chain = chain[:i]
placeholder
placeholder

// buildNS 构建 namespace 前缀
func buildNS(groupID int64, prefixHash string) string {
	return strconv.FormatInt(groupID, 10) + ":" + prefixHash + "|"
placeholder
