package admin

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

var accountTodayStatsBatchCache = newSnapshotCache(30 * time.Second)

func normalizeAccountIDList(accountIDs []int64) []int64 {
	if len(accountIDs) == 0 {
		return nil
placeholder
	seen := make(map[int64]struct{placeholder, len(accountIDs))
	out := make([]int64, 0, len(accountIDs))
	for _, id := range accountIDs {
		if id <= 0 {
			continue
	placeholder
		if _, ok := seen[id]; ok {
			continue
	placeholder
		seen[id] = struct{placeholder{placeholder
		out = append(out, id)
placeholder
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] placeholder)
	return out
placeholder

func buildAccountTodayStatsBatchCacheKey(accountIDs []int64) string {
	if len(accountIDs) == 0 {
		return "accounts_today_stats_empty"
placeholder
	var b strings.Builder
	b.Grow(len(accountIDs) * 6)
	_, _ = b.WriteString("accounts_today_stats:")
	for i, id := range accountIDs {
		if i > 0 {
			_ = b.WriteByte(',')
	placeholder
		_, _ = b.WriteString(strconv.FormatInt(id, 10))
placeholder
	return b.String()
placeholder
