package admin

import "sort"

func normalizeInt64IDList(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
placeholder

	out := make([]int64, 0, len(ids))
	seen := make(map[int64]struct{placeholder, len(ids))
	for _, id := range ids {
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
