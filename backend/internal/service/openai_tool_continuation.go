package service

import (
	"strings"

	"github.com/tidwall/gjson"
)

// ToolContinuationSignals 聚合工具续链相关信号，避免重复遍历 input。
type ToolContinuationSignals struct {
	HasFunctionCallOutput              bool
	HasFunctionCallOutputMissingCallID bool
	HasToolCallContext                 bool
	HasItemReference                   bool
	HasItemReferenceForAllCallIDs      bool
	FunctionCallOutputCallIDs          []string
placeholder

// FunctionCallOutputValidation 汇总 function_call_output 关联性校验结果。
type FunctionCallOutputValidation struct {
	HasFunctionCallOutput              bool
	HasToolCallContext                 bool
	HasFunctionCallOutputMissingCallID bool
	HasItemReferenceForAllCallIDs      bool
placeholder

func isCodexToolCallContextItemType(typ string) bool {
	switch strings.TrimSpace(typ) {
	case "tool_call",
		"function_call",
		"local_shell_call",
		"tool_search_call",
		"custom_tool_call",
		"mcp_tool_call":
		return true
	default:
		return false
placeholder
placeholder

func isCodexToolCallOutputItemType(typ string) bool {
	switch strings.TrimSpace(typ) {
	case "function_call_output",
		"tool_search_output",
		"custom_tool_call_output",
		"mcp_tool_call_output":
		return true
	default:
		return false
placeholder
placeholder

// NeedsToolContinuation 判定请求是否需要工具调用续链处理。
// 满足以下任一信号即视为续链：previous_response_id、input 内包含工具输出/item_reference、
// 或显式声明 tools/tool_choice。
func NeedsToolContinuation(reqBody map[string]any) bool {
	if reqBody == nil {
		return false
placeholder
	if hasNonEmptyString(reqBody["previous_response_id"]) {
		return true
placeholder
	if hasToolsSignal(reqBody) {
		return true
placeholder
	if hasToolChoiceSignal(reqBody) {
		return true
placeholder
	input, ok := reqBody["input"].([]any)
	if !ok {
		return false
placeholder
	for _, item := range input {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType, _ := itemMap["type"].(string)
		if isCodexToolCallItemType(itemType) || itemType == "item_reference" {
			return true
	placeholder
placeholder
	return false
placeholder

// AnalyzeToolContinuationSignals 单次遍历 input，提取工具输出/工具调用上下文/item_reference 相关信号。
// 字段名保留 FunctionCallOutput 是为了兼容既有调用点；语义覆盖 Codex 的所有工具输出
// （function_call_output/tool_search_output/custom_tool_call_output/mcp_tool_call_output）。
func AnalyzeToolContinuationSignals(reqBody map[string]any) ToolContinuationSignals {
	signals := ToolContinuationSignals{placeholder
	if reqBody == nil {
		return signals
placeholder
	input, ok := reqBody["input"].([]any)
	if !ok {
		return signals
placeholder

	var callIDs map[string]struct{placeholder
	var referenceIDs map[string]struct{placeholder

	for _, item := range input {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType, _ := itemMap["type"].(string)
		switch {
		case isCodexToolCallContextItemType(itemType):
			callID, _ := itemMap["call_id"].(string)
			if strings.TrimSpace(callID) != "" {
				signals.HasToolCallContext = true
		placeholder
		case isCodexToolCallOutputItemType(itemType):
			signals.HasFunctionCallOutput = true
			callID, _ := itemMap["call_id"].(string)
			callID = strings.TrimSpace(callID)
			if callID == "" {
				signals.HasFunctionCallOutputMissingCallID = true
				continue
		placeholder
			if callIDs == nil {
				callIDs = make(map[string]struct{placeholder)
		placeholder
			callIDs[callID] = struct{placeholder{placeholder
		case itemType == "item_reference":
			signals.HasItemReference = true
			idValue, _ := itemMap["id"].(string)
			idValue = strings.TrimSpace(idValue)
			if idValue == "" {
				continue
		placeholder
			if referenceIDs == nil {
				referenceIDs = make(map[string]struct{placeholder)
		placeholder
			referenceIDs[idValue] = struct{placeholder{placeholder
	placeholder
placeholder

	if len(callIDs) == 0 {
		return signals
placeholder
	signals.FunctionCallOutputCallIDs = make([]string, 0, len(callIDs))
	allReferenced := len(referenceIDs) > 0
	for callID := range callIDs {
		signals.FunctionCallOutputCallIDs = append(signals.FunctionCallOutputCallIDs, callID)
		if allReferenced {
			if _, ok := referenceIDs[callID]; !ok {
				allReferenced = false
		placeholder
	placeholder
placeholder
	signals.HasItemReferenceForAllCallIDs = allReferenced
	return signals
placeholder

// ValidateFunctionCallOutputContextBytes 基于 raw JSON 校验工具输出续链，避免 handler 预校验阶段全量解码大 input。
func ValidateFunctionCallOutputContextBytes(body []byte) FunctionCallOutputValidation {
	result := FunctionCallOutputValidation{placeholder
	if len(body) == 0 {
		return result
placeholder
	// handler 热路径只读扫描 input，避免 GetBytes 为大 Responses body 复制整段 JSON。
	input := parseRawJSONView(body).Get("input")
	if !input.IsArray() {
		return result
placeholder

	var callIDs map[string]struct{placeholder
	var referenceIDs map[string]struct{placeholder
	input.ForEach(func(_, item gjson.Result) bool {
		if !item.IsObject() {
			return true
	placeholder
		itemType := item.Get("type").String()
		switch {
		case isCodexToolCallOutputItemType(itemType):
			result.HasFunctionCallOutput = true
			callID := strings.TrimSpace(item.Get("call_id").String())
			if callID == "" {
				result.HasFunctionCallOutputMissingCallID = true
				return true
		placeholder
			if callIDs == nil {
				callIDs = make(map[string]struct{placeholder)
		placeholder
			callIDs[callID] = struct{placeholder{placeholder
		case isCodexToolCallContextItemType(itemType):
			if strings.TrimSpace(item.Get("call_id").String()) != "" {
				result.HasToolCallContext = true
		placeholder
		case itemType == "item_reference":
			idValue := strings.TrimSpace(item.Get("id").String())
			if idValue == "" {
				return true
		placeholder
			if referenceIDs == nil {
				referenceIDs = make(map[string]struct{placeholder)
		placeholder
			referenceIDs[idValue] = struct{placeholder{placeholder
	placeholder
		return !result.HasFunctionCallOutput || !result.HasToolCallContext
placeholder)
	if !result.HasFunctionCallOutput || result.HasToolCallContext || len(callIDs) == 0 || len(referenceIDs) == 0 {
		return result
placeholder
	allReferenced := true
	for callID := range callIDs {
		if _, ok := referenceIDs[callID]; !ok {
			allReferenced = false
			break
	placeholder
placeholder
	result.HasItemReferenceForAllCallIDs = allReferenced
	return result
placeholder

// ToolCallOutputContextCoverage 描述 input 中工具输出与可重建上下文的覆盖关系，
// 用于判断剥离 previous_response_id 后上游能否仅凭 input 重建工具续链。
type ToolCallOutputContextCoverage struct {
	HasFunctionCallOutput bool
	// ContextCoversAllCallIDs 表示每个工具输出的 call_id 都能在 input 内找到
	// 同 call_id 的工具调用上下文项或同 id 的 item_reference，且不存在缺失 call_id 的输出。
	// 任一输出无法由 input 自身重建时为 false，此时剥离 previous_response_id 会导致
	// 上游以 "No tool call found for function call output" 拒绝请求。
	ContextCoversAllCallIDs bool
placeholder

// AnalyzeToolCallOutputContextCoverageBytes 全量扫描 input，按 call_id 精确匹配工具输出
// 与可重建上下文。不能复用 ValidateFunctionCallOutputContextBytes 的 HasToolCallContext：
// 该标志只代表"存在某一个上下文项"，部分覆盖的续链仍会被上游拒绝。
func AnalyzeToolCallOutputContextCoverageBytes(body []byte) ToolCallOutputContextCoverage {
	coverage := ToolCallOutputContextCoverage{placeholder
	if len(body) == 0 {
		return coverage
placeholder
	input := parseRawJSONView(body).Get("input")
	if !input.IsArray() && !input.IsObject() {
		return coverage
placeholder

	missingCallID := false
	var outputCallIDs map[string]struct{placeholder
	var contextIDs map[string]struct{placeholder
	analyzeItem := func(item gjson.Result) {
		if !item.IsObject() {
			return
	placeholder
		itemType := item.Get("type").String()
		switch {
		case isCodexToolCallOutputItemType(itemType):
			coverage.HasFunctionCallOutput = true
			callID := strings.TrimSpace(item.Get("call_id").String())
			if callID == "" {
				missingCallID = true
				return
		placeholder
			if outputCallIDs == nil {
				outputCallIDs = make(map[string]struct{placeholder)
		placeholder
			outputCallIDs[callID] = struct{placeholder{placeholder
		case isCodexToolCallContextItemType(itemType):
			callID := strings.TrimSpace(item.Get("call_id").String())
			if callID == "" {
				return
		placeholder
			if contextIDs == nil {
				contextIDs = make(map[string]struct{placeholder)
		placeholder
			contextIDs[callID] = struct{placeholder{placeholder
		case itemType == "item_reference":
			idValue := strings.TrimSpace(item.Get("id").String())
			if idValue == "" {
				return
		placeholder
			if contextIDs == nil {
				contextIDs = make(map[string]struct{placeholder)
		placeholder
			contextIDs[idValue] = struct{placeholder{placeholder
	placeholder
placeholder
	if input.IsArray() {
		input.ForEach(func(_, item gjson.Result) bool {
			analyzeItem(item)
			return true
	placeholder)
placeholder else {
		analyzeItem(input)
placeholder

	if !coverage.HasFunctionCallOutput || missingCallID {
		return coverage
placeholder
	for callID := range outputCallIDs {
		if _, ok := contextIDs[callID]; !ok {
			return coverage
	placeholder
placeholder
	coverage.ContextCoversAllCallIDs = true
	return coverage
placeholder

// ValidateFunctionCallOutputContext 为 handler 提供低开销校验结果：
// 1) 无工具输出直接返回
// 2) 若已存在工具调用上下文则提前返回
// 3) 仅在无工具上下文时才构建 call_id / item_reference 集合
// 字段名保留 FunctionCallOutput 是为了兼容既有调用点；语义覆盖所有 Codex 工具输出。
func ValidateFunctionCallOutputContext(reqBody map[string]any) FunctionCallOutputValidation {
	result := FunctionCallOutputValidation{placeholder
	if reqBody == nil {
		return result
placeholder
	input, ok := reqBody["input"].([]any)
	if !ok {
		return result
placeholder

	for _, item := range input {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType, _ := itemMap["type"].(string)
		switch {
		case isCodexToolCallOutputItemType(itemType):
			result.HasFunctionCallOutput = true
		case isCodexToolCallContextItemType(itemType):
			callID, _ := itemMap["call_id"].(string)
			if strings.TrimSpace(callID) != "" {
				result.HasToolCallContext = true
		placeholder
	placeholder
		if result.HasFunctionCallOutput && result.HasToolCallContext {
			return result
	placeholder
placeholder

	if !result.HasFunctionCallOutput || result.HasToolCallContext {
		return result
placeholder

	callIDs := make(map[string]struct{placeholder)
	referenceIDs := make(map[string]struct{placeholder)
	for _, item := range input {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType, _ := itemMap["type"].(string)
		switch {
		case isCodexToolCallOutputItemType(itemType):
			callID, _ := itemMap["call_id"].(string)
			callID = strings.TrimSpace(callID)
			if callID == "" {
				result.HasFunctionCallOutputMissingCallID = true
				continue
		placeholder
			callIDs[callID] = struct{placeholder{placeholder
		case itemType == "item_reference":
			idValue, _ := itemMap["id"].(string)
			idValue = strings.TrimSpace(idValue)
			if idValue == "" {
				continue
		placeholder
			referenceIDs[idValue] = struct{placeholder{placeholder
	placeholder
placeholder

	if len(callIDs) == 0 || len(referenceIDs) == 0 {
		return result
placeholder
	allReferenced := true
	for callID := range callIDs {
		if _, ok := referenceIDs[callID]; !ok {
			allReferenced = false
			break
	placeholder
placeholder
	result.HasItemReferenceForAllCallIDs = allReferenced
	return result
placeholder

// HasFunctionCallOutput 判断 input 是否包含任意 Codex 工具输出，用于触发续链校验。
// 名称保留 function_call_output 是为了兼容既有调用点。
func HasFunctionCallOutput(reqBody map[string]any) bool {
	return AnalyzeToolContinuationSignals(reqBody).HasFunctionCallOutput
placeholder

// HasToolCallContext 判断 input 是否包含带 call_id 的工具调用上下文，
// 用于判断工具输出是否具备可关联的上下文。
func HasToolCallContext(reqBody map[string]any) bool {
	return AnalyzeToolContinuationSignals(reqBody).HasToolCallContext
placeholder

// FunctionCallOutputCallIDs 提取 input 中工具输出的 call_id 集合。
// 仅返回非空 call_id，用于与 item_reference.id 做匹配校验。
func FunctionCallOutputCallIDs(reqBody map[string]any) []string {
	return AnalyzeToolContinuationSignals(reqBody).FunctionCallOutputCallIDs
placeholder

// HasFunctionCallOutputMissingCallID 判断是否存在缺少 call_id 的工具输出。
func HasFunctionCallOutputMissingCallID(reqBody map[string]any) bool {
	return AnalyzeToolContinuationSignals(reqBody).HasFunctionCallOutputMissingCallID
placeholder

// HasItemReferenceForCallIDs 判断 item_reference.id 是否覆盖所有 call_id。
// 用于仅依赖引用项完成续链场景的校验。
func HasItemReferenceForCallIDs(reqBody map[string]any, callIDs []string) bool {
	if reqBody == nil || len(callIDs) == 0 {
		return false
placeholder
	input, ok := reqBody["input"].([]any)
	if !ok {
		return false
placeholder
	referenceIDs := make(map[string]struct{placeholder)
	for _, item := range input {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType, _ := itemMap["type"].(string)
		if itemType != "item_reference" {
			continue
	placeholder
		idValue, _ := itemMap["id"].(string)
		idValue = strings.TrimSpace(idValue)
		if idValue == "" {
			continue
	placeholder
		referenceIDs[idValue] = struct{placeholder{placeholder
placeholder
	if len(referenceIDs) == 0 {
		return false
placeholder
	for _, callID := range callIDs {
		if _, ok := referenceIDs[strings.TrimSpace(callID)]; !ok {
			return false
	placeholder
placeholder
	return true
placeholder

// hasNonEmptyString 判断字段是否为非空字符串。
func hasNonEmptyString(value any) bool {
	stringValue, ok := value.(string)
	return ok && strings.TrimSpace(stringValue) != ""
placeholder

// hasToolsSignal 判断 tools 字段是否显式声明（存在且不为空）。
func hasToolsSignal(reqBody map[string]any) bool {
	raw, exists := reqBody["tools"]
	if !exists || raw == nil {
		return false
placeholder
	if tools, ok := raw.([]any); ok {
		return len(tools) > 0
placeholder
	return false
placeholder

// hasToolChoiceSignal 判断 tool_choice 是否显式声明（非空或非 nil）。
func hasToolChoiceSignal(reqBody map[string]any) bool {
	raw, exists := reqBody["tool_choice"]
	if !exists || raw == nil {
		return false
placeholder
	switch value := raw.(type) {
	case string:
		return strings.TrimSpace(value) != ""
	case map[string]any:
		return len(value) > 0
	default:
		return false
placeholder
placeholder
