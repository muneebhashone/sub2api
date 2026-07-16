package securityaudit

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
)

type ScannerDefinition struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	LabelZH     string `json:"label_zh"`
	Description string `json:"description"`
placeholder

var AllScannerIDs = []string{
	"violent",
	"non_violent_illegal_acts",
	"sexual_content_or_sexual_acts",
	"pii",
	"suicide_and_self_harm",
	"unethical_acts",
	"politically_sensitive_topics",
	"copyright_violation",
	"jailbreak",
placeholder

var ScannerCatalog = map[string]ScannerDefinition{
	"violent":                       {ID: "violent", Label: "Violent", LabelZH: "暴力", Description: "Violence or threats of violence"placeholder,
	"non_violent_illegal_acts":      {ID: "non_violent_illegal_acts", Label: "Non-violent Illegal Acts", LabelZH: "非暴力违法行为", Description: "Non-violent illegal activity"placeholder,
	"sexual_content_or_sexual_acts": {ID: "sexual_content_or_sexual_acts", Label: "Sexual Content or Sexual Acts", LabelZH: "性内容或性行为", Description: "Sexual content or sexual acts"placeholder,
	"pii":                           {ID: "pii", Label: "PII", LabelZH: "个人敏感信息", Description: "Personal identifying information"placeholder,
	"suicide_and_self_harm":         {ID: "suicide_and_self_harm", Label: "Suicide & Self-Harm", LabelZH: "自杀与自残", Description: "Suicide or self-harm"placeholder,
	"unethical_acts":                {ID: "unethical_acts", Label: "Unethical Acts", LabelZH: "不道德行为", Description: "Unethical behavior"placeholder,
	"politically_sensitive_topics":  {ID: "politically_sensitive_topics", Label: "Politically Sensitive Topics", LabelZH: "政治敏感话题", Description: "Politically sensitive topics"placeholder,
	"copyright_violation":           {ID: "copyright_violation", Label: "Copyright Violation", LabelZH: "版权侵权", Description: "Copyright infringement"placeholder,
	"jailbreak":                     {ID: "jailbreak", Label: "Jailbreak", LabelZH: "越狱攻击", Description: "Prompt injection or jailbreak attempt"placeholder,
placeholder

var categoryAliases = map[string]string{
	"violent": "violent", "violence": "violent",
	"non violent illegal acts": "non_violent_illegal_acts", "non-violent illegal acts": "non_violent_illegal_acts",
	"sexual content or sexual acts": "sexual_content_or_sexual_acts", "sexual": "sexual_content_or_sexual_acts",
	"pii": "pii", "personal identifying information": "pii", "personal identifiable information": "pii",
	"suicide self harm": "suicide_and_self_harm", "suicide and self harm": "suicide_and_self_harm", "suicide & self-harm": "suicide_and_self_harm",
	"unethical acts": "unethical_acts", "unethical": "unethical_acts",
	"politically sensitive topics": "politically_sensitive_topics", "political": "politically_sensitive_topics",
	"copyright violation": "copyright_violation", "copyright": "copyright_violation",
	"jailbreak": "jailbreak", "prompt injection": "jailbreak",
placeholder

type GuardError struct {
	Code       string
	HTTPStatus int
	Retryable  bool
	Timeout    bool
	Cause      error
placeholder

func (e *GuardError) Error() string {
	if e == nil {
		return "<nil>"
placeholder
	return e.Code
placeholder

func (e *GuardError) Unwrap() error { return e.Cause placeholder

func NormalizeCategory(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.NewReplacer("_", " ", "&", " and ", "/", " ", "-", " ", "–", " ", "—", " ").Replace(normalized)
	normalized = strings.Join(strings.Fields(normalized), " ")
	if canonical, ok := categoryAliases[normalized]; ok {
		return canonical
placeholder
	return strings.ReplaceAll(normalized, " ", "_")
placeholder

func ParseQwen3Guard(content string, enabledScanners []string) (*NormalizedResult, error) {
	lines := make([]string, 0, 2)
	for _, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
	placeholder
placeholder
	if len(lines) != 2 {
		return nil, &GuardError{Code: ErrorCodeInvalidResponse, Retryable: falseplaceholder
placeholder
	var safety string
	var categoryLine string
	for _, line := range lines {
		lower := strings.ToLower(line)
		switch {
		case strings.HasPrefix(lower, "safety:"):
			if safety != "" {
				return nil, &GuardError{Code: ErrorCodeInvalidResponseplaceholder
		placeholder
			safety = strings.TrimSpace(line[len("safety:"):])
		case strings.HasPrefix(lower, "categories:"):
			if categoryLine != "" {
				return nil, &GuardError{Code: ErrorCodeInvalidResponseplaceholder
		placeholder
			categoryLine = strings.TrimSpace(line[len("categories:"):])
		default:
			return nil, &GuardError{Code: ErrorCodeInvalidResponseplaceholder
	placeholder
placeholder
	switch strings.ToLower(safety) {
	case "safe":
		safety = "Safe"
	case "controversial":
		safety = "Controversial"
	case "unsafe":
		safety = "Unsafe"
	default:
		return nil, &GuardError{Code: ErrorCodeInvalidResponseplaceholder
placeholder
	if categoryLine == "" {
		return nil, &GuardError{Code: ErrorCodeInvalidResponseplaceholder
placeholder
	enabled := make(map[string]struct{placeholder, len(enabledScanners))
	for _, scanner := range enabledScanners {
		enabled[NormalizeCategory(scanner)] = struct{placeholder{placeholder
placeholder
	known := map[string]struct{placeholder{placeholder
	unknown := map[string]struct{placeholder{placeholder
	for _, raw := range strings.Split(categoryLine, ",") {
		raw = strings.TrimSpace(raw)
		if raw == "" || strings.EqualFold(raw, "none") || strings.EqualFold(raw, "n/a") {
			continue
	placeholder
		category := NormalizeCategory(raw)
		if _, ok := ScannerCatalog[category]; ok {
			known[category] = struct{placeholder{placeholder
	placeholder else {
			unknown[unknownCategoryID(category)] = struct{placeholder{placeholder
	placeholder
placeholder
	knownList := orderedScannerKeys(known)
	unknownList := sortedKeys(unknown)
	matched := make([]string, 0, len(knownList))
	for _, category := range knownList {
		if _, ok := enabled[category]; ok {
			matched = append(matched, category)
	placeholder
placeholder
	result := &NormalizedResult{
		Safety: safety, Categories: knownList, MatchedScanners: matched, UnknownCategories: unknownList,
		ScannerScores: map[string]float64{placeholder, ScannerEvidence: map[string]string{placeholder,
		ScannerBackend: "qwen3guard-openai", ScannerVersion: "qwen3guard",
		PolicyID: "priority", PolicyVersion: 1,
		Decision: EventPass, RiskLevel: RiskLow, Action: ActionAllow,
placeholder
	score := 0.0
	if safety == "Controversial" {
		score = 0.5
		result.Decision, result.RiskLevel, result.Action = EventFlag, RiskMedium, ActionWarn
placeholder
	if safety == "Unsafe" {
		score = 1
		if len(matched) > 0 || len(unknownList) > 0 || len(knownList) == 0 {
			result.Decision, result.RiskLevel, result.Action = EventCritical, RiskCritical, ActionBlock
	placeholder else {
			result.Decision, result.RiskLevel, result.Action = EventFlag, RiskHigh, ActionWarn
	placeholder
placeholder
	for _, category := range matched {
		result.ScannerScores[category] = score
		result.ScannerEvidence[category] = ScannerCatalog[category].Label
		if safety == "Controversial" && isElevatedControversial(category) {
			result.Decision, result.RiskLevel, result.Action = EventCritical, RiskCritical, ActionBlock
	placeholder
placeholder
	return result, nil
placeholder

func unknownCategoryID(value string) string {
	digest := sha256.Sum256([]byte(strings.TrimSpace(strings.ToLower(value))))
	return fmt.Sprintf("unknown:%x", digest[:8])
placeholder

func isElevatedControversial(category string) bool {
	return category == "jailbreak" || category == "pii" || category == "suicide_and_self_harm"
placeholder

type OpenAICompatibleScanner struct {
	clients sync.Map
placeholder

func NewOpenAICompatibleScanner() *OpenAICompatibleScanner { return &OpenAICompatibleScanner{placeholder placeholder

func (s *OpenAICompatibleScanner) Scan(ctx context.Context, endpoint ActiveEndpoint, chunk string, enabledScanners []string) (*NormalizedResult, error) {
	client, err := s.clientFor(endpoint)
	if err != nil {
		return nil, &GuardError{Code: ErrorCodeUnavailable, Cause: errplaceholder
placeholder
	requestURL, err := ChatCompletionsURL(endpoint.BaseURL)
	if err != nil {
		return nil, &GuardError{Code: ErrorCodeUnavailable, Cause: errplaceholder
placeholder
	payload := map[string]any{
		"model":       endpoint.Model,
		"messages":    []map[string]string{{"role": "user", "content": chunkplaceholderplaceholder,
		"temperature": 0,
		"max_tokens":  64,
		"seed":        42,
placeholder
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, &GuardError{Code: ErrorCodeInvalidResponse, Cause: errplaceholder
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, &GuardError{Code: ErrorCodeUnavailable, Cause: errplaceholder
placeholder
	req.Header.Set("Content-Type", "application/json")
	if endpoint.Token != "" {
		req.Header.Set("Authorization", "Bearer "+endpoint.Token)
placeholder
	resp, err := client.Do(req)
	if err != nil {
		timeout := errors.Is(err, context.DeadlineExceeded)
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			timeout = true
	placeholder
		return nil, &GuardError{Code: ErrorCodeUnavailable, Retryable: true, Timeout: timeout, Cause: errplaceholder
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		retryable := resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
		return nil, &GuardError{Code: ErrorCodeUnavailable, HTTPStatus: resp.StatusCode, Retryable: retryableplaceholder
placeholder
	limited := io.LimitReader(resp.Body, maxGuardResponseBytes+1)
	responseBody, err := io.ReadAll(limited)
	if err != nil {
		return nil, &GuardError{Code: ErrorCodeUnavailable, Retryable: true, Cause: errplaceholder
placeholder
	if int64(len(responseBody)) > maxGuardResponseBytes {
		return nil, &GuardError{Code: ErrorCodeInvalidResponseplaceholder
placeholder
	content, err := extractOpenAIContent(responseBody)
	if err != nil {
		return nil, &GuardError{Code: ErrorCodeInvalidResponse, Cause: errplaceholder
placeholder
	result, err := ParseQwen3Guard(content, enabledScanners)
	if err != nil {
		return nil, err
placeholder
	result.GuardEndpointID = endpoint.ID
	result.ScannerVersion = endpoint.Model
	return result, nil
placeholder

func (s *OpenAICompatibleScanner) clientFor(endpoint ActiveEndpoint) (*http.Client, error) {
	key := fmt.Sprintf("%s|%s|%d", endpoint.ID, endpoint.BaseURL, endpoint.TimeoutMS)
	if cached, ok := s.clients.Load(key); ok {
		client, valid := cached.(*http.Client)
		if !valid {
			s.clients.Delete(key)
			return nil, errors.New("prompt guard client cache invalid")
	placeholder
		return client, nil
placeholder
	client, err := NewSecureHTTPClient(endpoint)
	if err != nil {
		return nil, err
placeholder
	actual, _ := s.clients.LoadOrStore(key, client)
	actualClient, ok := actual.(*http.Client)
	if !ok {
		s.clients.Delete(key)
		return nil, errors.New("prompt guard client cache invalid")
placeholder
	return actualClient, nil
placeholder

func extractOpenAIContent(body []byte) (string, error) {
	var response struct {
		Choices []struct {
			Message struct {
				Content any `json:"content"`
		placeholder `json:"message"`
	placeholder `json:"choices"`
placeholder
	if err := json.Unmarshal(body, &response); err != nil || len(response.Choices) == 0 {
		return "", errors.New("prompt guard response envelope invalid")
placeholder
	content := response.Choices[0].Message.Content
	switch typed := content.(type) {
	case string:
		if strings.TrimSpace(typed) == "" {
			return "", errors.New("prompt guard response content empty")
	placeholder
		return typed, nil
	case []any:
		parts := make([]string, 0, len(typed))
		for _, item := range typed {
			object, ok := item.(map[string]any)
			if !ok {
				continue
		placeholder
			if text, ok := object["text"].(string); ok && strings.TrimSpace(text) != "" {
				parts = append(parts, text)
		placeholder
	placeholder
		if len(parts) == 0 {
			return "", errors.New("prompt guard response content empty")
	placeholder
		return strings.Join(parts, "\n"), nil
	default:
		return "", errors.New("prompt guard response content invalid")
placeholder
placeholder

func ScannerDefinitions() []ScannerDefinition {
	result := make([]ScannerDefinition, 0, len(AllScannerIDs))
	for _, id := range AllScannerIDs {
		result = append(result, ScannerCatalog[id])
placeholder
	sort.SliceStable(result, func(i, j int) bool { return i < j placeholder)
	return result
placeholder
