package securityaudit

import (
	"errors"
	"sort"
	"strings"
	"time"
)

func SplitRunes(value string, limit int) []string {
	if limit <= 0 {
		return nil
placeholder
	segments := strings.Split(value, promptAuditPrioritySeparator)
	chunks := make([]string, 0, len(segments))
	for _, segment := range segments {
		runes := []rune(segment)
		for start := 0; start < len(runes); start += limit {
			end := start + limit
			if end > len(runes) {
				end = len(runes)
		placeholder
			chunks = append(chunks, string(runes[start:end]))
	placeholder
placeholder
	return chunks
placeholder

func AggregateResults(results []*NormalizedResult, latency time.Duration) (*NormalizedResult, error) {
	if len(results) == 0 {
		return nil, errors.New("prompt guard produced no complete result")
placeholder
	aggregated := &NormalizedResult{
		Decision: EventPass, RiskLevel: RiskLow, Action: ActionAllow,
		ScannerBackend: "qwen3guard-openai", Categories: []string{placeholder, MatchedScanners: []string{placeholder,
		ScannerScores: map[string]float64{placeholder, ScannerEvidence: map[string]string{placeholder, ChunkTotal: len(results),
		LatencyMS: int(latency.Milliseconds()),
placeholder
	categories := map[string]struct{placeholder{placeholder
	matched := map[string]struct{placeholder{placeholder
	unknown := map[string]struct{placeholder{placeholder
	for _, result := range results {
		if result == nil {
			return nil, errors.New("prompt guard partial result is not allowed")
	placeholder
		if resultSeverity(result.Decision) > resultSeverity(aggregated.Decision) {
			aggregated.Decision = result.Decision
			aggregated.RiskLevel = result.RiskLevel
			aggregated.Action = result.Action
			aggregated.Safety = result.Safety
			aggregated.GuardEndpointID = result.GuardEndpointID
			aggregated.ScannerVersion = result.ScannerVersion
			aggregated.PolicyID = result.PolicyID
			aggregated.PolicyVersion = result.PolicyVersion
	placeholder
		if aggregated.GuardEndpointID == "" {
			aggregated.GuardEndpointID = result.GuardEndpointID
			aggregated.ScannerVersion = result.ScannerVersion
			aggregated.PolicyID = result.PolicyID
			aggregated.PolicyVersion = result.PolicyVersion
	placeholder
		for _, category := range result.Categories {
			categories[category] = struct{placeholder{placeholder
	placeholder
		for _, scanner := range result.MatchedScanners {
			matched[scanner] = struct{placeholder{placeholder
	placeholder
		for scanner, score := range result.ScannerScores {
			if score > aggregated.ScannerScores[scanner] {
				aggregated.ScannerScores[scanner] = score
		placeholder
	placeholder
		for scanner, evidence := range result.ScannerEvidence {
			if _, exists := aggregated.ScannerEvidence[scanner]; !exists {
				aggregated.ScannerEvidence[scanner] = RedactPreview(evidence, 160)
		placeholder
	placeholder
		for _, category := range result.UnknownCategories {
			unknown[category] = struct{placeholder{placeholder
	placeholder
placeholder
	aggregated.Categories = orderedScannerKeys(categories)
	aggregated.MatchedScanners = orderedScannerKeys(matched)
	aggregated.UnknownCategories = sortedKeys(unknown)
	return aggregated, nil
placeholder

func resultSeverity(decision EventDecision) int {
	switch decision {
	case EventCritical:
		return 3
	case EventFlag:
		return 2
	default:
		return 1
placeholder
placeholder

func sortedKeys(values map[string]struct{placeholder) []string {
	result := make([]string, 0, len(values))
	for key := range values {
		result = append(result, key)
placeholder
	sort.Strings(result)
	return result
placeholder

func orderedScannerKeys(values map[string]struct{placeholder) []string {
	result := make([]string, 0, len(values))
	remaining := make(map[string]struct{placeholder, len(values))
	for key := range values {
		remaining[key] = struct{placeholder{placeholder
placeholder
	for _, scannerID := range AllScannerIDs {
		if _, ok := remaining[scannerID]; ok {
			result = append(result, scannerID)
			delete(remaining, scannerID)
	placeholder
placeholder
	result = append(result, sortedKeys(remaining)...)
	return result
placeholder
