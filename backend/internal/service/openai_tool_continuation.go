package service

import "strings"

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

// NeedsToolContinuation 判定请求是否需要工具调用续链处理。
// 满足以下任一信号即视为续链：previous_response_id、input 内包含 function_call_output/item_reference、
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
		if itemType == "function_call_output" || itemType == "item_reference" {
			return true
	placeholder
placeholder
	return false
placeholder

// AnalyzeToolContinuationSignals 单次遍历 input，提取 function_call_output/tool_call/item_reference 相关信号。
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
		switch itemType {
		case "tool_call", "function_call":
			callID, _ := itemMap["call_id"].(string)
			if strings.TrimSpace(callID) != "" {
				signals.HasToolCallContext = true
		placeholder
		case "function_call_output":
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
		case "item_reference":
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

// ValidateFunctionCallOutputContext 为 handler 提供低开销校验结果：
// 1) 无 function_call_output 直接返回
// 2) 若已存在 tool_call/function_call 上下文则提前返回
// 3) 仅在无工具上下文时才构建 call_id / item_reference 集合
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
		switch itemType {
		case "function_call_output":
			result.HasFunctionCallOutput = true
		case "tool_call", "function_call":
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
		switch itemType {
		case "function_call_output":
			callID, _ := itemMap["call_id"].(string)
			callID = strings.TrimSpace(callID)
			if callID == "" {
				result.HasFunctionCallOutputMissingCallID = true
				continue
		placeholder
			callIDs[callID] = struct{placeholder{placeholder
		case "item_reference":
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

// HasFunctionCallOutput 判断 input 是否包含 function_call_output，用于触发续链校验。
func HasFunctionCallOutput(reqBody map[string]any) bool {
	return AnalyzeToolContinuationSignals(reqBody).HasFunctionCallOutput
placeholder

// HasToolCallContext 判断 input 是否包含带 call_id 的 tool_call/function_call，
// 用于判断 function_call_output 是否具备可关联的上下文。
func HasToolCallContext(reqBody map[string]any) bool {
	return AnalyzeToolContinuationSignals(reqBody).HasToolCallContext
placeholder

// FunctionCallOutputCallIDs 提取 input 中 function_call_output 的 call_id 集合。
// 仅返回非空 call_id，用于与 item_reference.id 做匹配校验。
func FunctionCallOutputCallIDs(reqBody map[string]any) []string {
	return AnalyzeToolContinuationSignals(reqBody).FunctionCallOutputCallIDs
placeholder

// HasFunctionCallOutputMissingCallID 判断是否存在缺少 call_id 的 function_call_output。
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
