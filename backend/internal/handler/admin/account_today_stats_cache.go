package admin

import (
	"strconv"
	"strings"
	"time"
)

var accountTodayStatsBatchCache = newSnapshotCache(30 * time.Second)

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
