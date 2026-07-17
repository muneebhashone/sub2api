package securityaudit

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseQwen3GuardStrictAndPolicy(t *testing.T) {
	tests := []struct {
		name, output string
		enabled      []string
		decision     EventDecision
		action       Action
		wantErr      bool
placeholder{
		{"safe", "Safety: Safe\nCategories: None", AllScannerIDs, EventPass, ActionAllow, falseplaceholder,
		{"controversial", "Safety: Controversial\nCategories: Violent", AllScannerIDs, EventFlag, ActionWarn, falseplaceholder,
		{"controversial pii escalates", "Safety: Controversial\nCategories: PII", AllScannerIDs, EventCritical, ActionBlock, falseplaceholder,
		{"unsafe", "Safety: Unsafe\nCategories: Jailbreak", AllScannerIDs, EventCritical, ActionBlock, falseplaceholder,
		{"unknown unsafe", "Safety: Unsafe\nCategories: Future Risk", AllScannerIDs, EventCritical, ActionBlock, falseplaceholder,
		{"disabled unsafe warns", "Safety: Unsafe\nCategories: Violent", []string{"PII"placeholder, EventFlag, ActionWarn, falseplaceholder,
		{"extra explanation", "Safety: Safe\nCategories: None\nThis is safe", AllScannerIDs, "", "", trueplaceholder,
		{"duplicate", "Safety: Safe\nSafety: Safe", AllScannerIDs, "", "", trueplaceholder,
		{"duplicate categories", "Safety: Safe\nCategories: None\nCategories: PII", AllScannerIDs, "", "", trueplaceholder,
		{"missing categories", "Safety: Safe\n", AllScannerIDs, "", "", trueplaceholder,
		{"unknown safety", "Safety: Maybe\nCategories: PII", AllScannerIDs, "", "", trueplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseQwen3Guard(tt.output, tt.enabled)
			if tt.wantErr {
			placeholder
				return
		placeholder
		placeholder
			require.Equal(t, tt.decision, result.Decision)
			require.Equal(t, tt.action, result.Action)
	placeholder)
placeholder
placeholder

func TestQwen3GuardOfficialCategoriesAliasesAndUnknownAreStable(t *testing.T) {
	official := "Violent, Non-violent Illegal Acts, Sexual Content or Sexual Acts, PII, Suicide & Self-Harm, Unethical Acts, Politically Sensitive Topics, Copyright Violation, Jailbreak"
	result, err := ParseQwen3Guard("Safety: Unsafe\nCategories: "+official, AllScannerIDs)
placeholder
	require.Equal(t, AllScannerIDs, result.MatchedScanners)
	require.Empty(t, result.UnknownCategories)
	require.Equal(t, "priority", result.PolicyID)
	require.Equal(t, 1, result.PolicyVersion)

	aliases := map[string]string{
		"violence": "violent", "non_violent_illegal_acts": "non_violent_illegal_acts",
		"sexual": "sexual_content_or_sexual_acts", "personal identifiable information": "pii",
		"suicide/self harm": "suicide_and_self_harm", "unethical": "unethical_acts",
		"political": "politically_sensitive_topics", "copyright": "copyright_violation",
		"prompt injection": "jailbreak",
placeholder
	for alias, canonical := range aliases {
		require.Equal(t, canonical, NormalizeCategory(alias), alias)
placeholder

	const canary = "PROMPT_CANARY_RAW_UNKNOWN_CATEGORY"
	unknown, err := ParseQwen3Guard("Safety: Unsafe\nCategories: "+canary, AllScannerIDs)
placeholder
	require.Len(t, unknown.UnknownCategories, 1)
	require.NotContains(t, unknown.UnknownCategories[0], "canary")
	require.NotContains(t, unknown.UnknownCategories[0], "raw")
	require.Contains(t, unknown.UnknownCategories[0], "unknown:")
placeholder

func TestExtractOpenAIContentSupportsStringAndTextBlocks(t *testing.T) {
	content, err := extractOpenAIContent([]byte(`{"choices":[{"message":{"content":"Safety: Safe\nCategories: None"placeholderplaceholder]placeholder`))
placeholder
	require.Equal(t, "Safety: Safe\nCategories: None", content)
	content, err = extractOpenAIContent([]byte(`{"choices":[{"message":{"content":[{"type":"text","text":"Safety: Safe"placeholder,{"type":"text","text":"Categories: None"placeholder]placeholderplaceholder]placeholder`))
placeholder
	require.Equal(t, "Safety: Safe\nCategories: None", content)
	for _, body := range []string{`{placeholder`, `{"choices":[]placeholder`, `{"choices":[{"message":{"content":nullplaceholderplaceholder]placeholder`placeholder {
		_, err := extractOpenAIContent([]byte(body))
	placeholder
placeholder
placeholder

func TestAggregateRequiresEveryResult(t *testing.T) {
	_, err := AggregateResults([]*NormalizedResult{{Decision: EventPass, Action: ActionAllowplaceholder, nilplaceholder, 0)
placeholder
	result, err := AggregateResults([]*NormalizedResult{
		{Decision: EventPass, RiskLevel: RiskLow, Action: ActionAllow, Categories: []string{"pii"placeholderplaceholder,
		{Decision: EventCritical, RiskLevel: RiskCritical, Action: ActionBlock, Categories: []string{"jailbreak"placeholderplaceholder,
placeholder, 0)
placeholder
	require.Equal(t, EventCritical, result.Decision)
	require.Equal(t, ActionBlock, result.Action)
	require.Equal(t, []string{"pii", "jailbreak"placeholder, result.Categories)
placeholder

func TestAggregateDeduplicatesFactsAndUsesMostSevereEndpointMetadata(t *testing.T) {
	result, err := AggregateResults([]*NormalizedResult{
		{Decision: EventPass, RiskLevel: RiskLow, Action: ActionAllow, Safety: "Safe", Categories: []string{"pii"placeholder, MatchedScanners: []string{"pii"placeholder, ScannerScores: map[string]float64{"pii": 0placeholder, ScannerEvidence: map[string]string{"pii": "first"placeholder, GuardEndpointID: "safe-node", ScannerVersion: "safe-version", PolicyID: "priority", PolicyVersion: 1placeholder,
		{Decision: EventCritical, RiskLevel: RiskCritical, Action: ActionBlock, Safety: "Unsafe", Categories: []string{"pii", "jailbreak"placeholder, MatchedScanners: []string{"pii", "jailbreak"placeholder, ScannerScores: map[string]float64{"pii": 1, "jailbreak": 1placeholder, ScannerEvidence: map[string]string{"pii": "second", "jailbreak": "blocked"placeholder, GuardEndpointID: "block-node", ScannerVersion: "block-version", PolicyID: "priority", PolicyVersion: 2placeholder,
placeholder, 7*time.Millisecond)
placeholder
	require.Equal(t, []string{"pii", "jailbreak"placeholder, result.Categories)
	require.Equal(t, []string{"pii", "jailbreak"placeholder, result.MatchedScanners)
	require.Equal(t, "first", result.ScannerEvidence["pii"], "evidence is deterministically first-seen")
	require.Equal(t, "block-node", result.GuardEndpointID)
	require.Equal(t, "block-version", result.ScannerVersion)
	require.Equal(t, 2, result.PolicyVersion)
	require.Equal(t, 7, result.LatencyMS)
placeholder

func TestIssueSummariesAreDeterministicRedactedDerivedDTOs(t *testing.T) {
	const canary = "PROMPT_CANARY_EVIDENCE_SECRET"
	result := NormalizedResult{
		Decision: EventCritical, RiskLevel: RiskCritical, Action: ActionBlock,
		Categories: []string{"jailbreak", "pii"placeholder, MatchedScanners: []string{"pii"placeholder,
		ScannerScores: map[string]float64{"pii": 1placeholder, ScannerEvidence: map[string]string{"pii": canaryplaceholder,
		UnknownCategories: []string{unknownCategoryID("future risk")placeholder,
placeholder
	summaries := BuildIssueSummaries(result)
	require.Len(t, summaries, 3, "known categories are not hidden merely because policy disabled one")
	raw, err := json.Marshal(summaries)
placeholder
	require.NotContains(t, string(raw), canary)
	for _, summary := range summaries {
		require.NotEmpty(t, summary.Title)
		require.NotEmpty(t, summary.Description)
		require.NotEmpty(t, summary.Code)
		require.NotEmpty(t, summary.EvidenceHash)
placeholder
placeholder
