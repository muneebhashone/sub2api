package service

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const openAIResponsesNamespaceNamesContextKey = "openai_responses_namespace_names"

// shouldFlattenOpenAIResponsesNamespaces 判定原生 Responses 转发前是否摊平
// Codex namespace 工具。WSv2 上游原生支持 namespace，且 WS 出口
// （openai_ws_forwarder_v2）原样转发上游事件、不经 HTTP 回程还原，摊平后的
// 平名无法还原会破坏客户端工具匹配，因此实际走 WSv2 分支的请求保持 namespace
// 原样。透传账号先于 WSv2 分支经 HTTP 转发返回，仍需摊平。
func shouldFlattenOpenAIResponsesNamespaces(account *Account, transport OpenAIUpstreamTransport, passthroughEnabled bool) bool {
	if account == nil || !account.IsOpenAIOAuth() {
		return false
placeholder
	if transport == OpenAIUpstreamTransportResponsesWebsocketV2 && !passthroughEnabled {
		return false
placeholder
	return true
placeholder

// shouldStripOpenAIResponsesInputNamespaces removes residual input item
// namespaces for OpenAI OAuth and API Key HTTP forwarding. Native WSv2 keeps
// namespaces because that protocol supports them and does not restore payloads.
func shouldStripOpenAIResponsesInputNamespaces(account *Account, transport OpenAIUpstreamTransport, passthroughEnabled bool) bool {
	if account == nil || (!account.IsOpenAIOAuth() && !account.IsOpenAIApiKey()) {
		return false
placeholder
	if transport == OpenAIUpstreamTransportResponsesWebsocketV2 && !passthroughEnabled {
		return false
placeholder
	return true
placeholder

func flattenOpenAIResponsesNamespaces(c *gin.Context, body []byte) ([]byte, error) {
	if !bytes.Contains(body, []byte(`"namespace"`)) {
		return body, nil
placeholder
	var requestBody map[string]any
	if err := json.Unmarshal(body, &requestBody); err != nil {
		return body, fmt.Errorf("decode OpenAI namespace body: %w", err)
placeholder
	names, changed, err := apicompat.FlattenResponsesNamespacesExcept(requestBody, map[string]bool{"image_gen": trueplaceholder)
	if err != nil {
		return body, err
placeholder
	if !changed {
		return body, nil
placeholder
	rebuilt, err := marshalOpenAIUpstreamJSON(requestBody)
	if err != nil {
		return body, fmt.Errorf("encode OpenAI namespace body: %w", err)
placeholder
	setOpenAIResponsesNamespaceNames(c, names)
	return rebuilt, nil
placeholder

// stripOpenAIResponsesInputNamespaces removes namespace only from direct input
// array items. Namespace declarations and nested namespace fields are left
// untouched. Rebuilding the input array once keeps this linear for long
// histories and avoids decoding JSON numbers through float64.
func stripOpenAIResponsesInputNamespaces(body []byte) ([]byte, error) {
	if !bytes.Contains(body, []byte(`"namespace"`)) {
		return body, nil
placeholder
	input := gjson.GetBytes(body, "input")
	if !input.IsArray() {
		return body, nil
placeholder

	var rebuilt bytes.Buffer
	rebuilt.Grow(len(input.Raw))
	_ = rebuilt.WriteByte('[')
	changed := false
	first := true
	var stripErr error
	input.ForEach(func(_, item gjson.Result) bool {
		if !first {
			_ = rebuilt.WriteByte(',')
	placeholder
		first = false
		itemBody := []byte(item.Raw)
		if item.IsObject() && item.Get("namespace").Exists() {
			itemBody, stripErr = sjson.DeleteBytes(itemBody, "namespace")
			if stripErr != nil {
				return false
		placeholder
			changed = true
	placeholder
		_, _ = rebuilt.Write(itemBody)
		return true
placeholder)
	if stripErr != nil {
		return body, fmt.Errorf("delete OpenAI input namespace: %w", stripErr)
placeholder
	if !changed {
		return body, nil
placeholder
	_ = rebuilt.WriteByte(']')
	stripped, err := sjson.SetRawBytes(body, "input", rebuilt.Bytes())
	if err != nil {
		return body, fmt.Errorf("replace OpenAI input after namespace deletion: %w", err)
placeholder
	return stripped, nil
placeholder

func setOpenAIResponsesNamespaceNames(c *gin.Context, names map[string]apicompat.ResponsesNamespaceName) {
	if c != nil && len(names) > 0 {
		c.Set(openAIResponsesNamespaceNamesContextKey, names)
placeholder
placeholder

func openAIResponsesNamespaceNames(c *gin.Context) map[string]apicompat.ResponsesNamespaceName {
	if c == nil {
		return nil
placeholder
	value, ok := c.Get(openAIResponsesNamespaceNamesContextKey)
	if !ok {
		return nil
placeholder
	names, _ := value.(map[string]apicompat.ResponsesNamespaceName)
	return names
placeholder

func restoreOpenAIResponsesNamespacePayload(c *gin.Context, payload []byte) ([]byte, error) {
	names := openAIResponsesNamespaceNames(c)
	if len(names) == 0 || !json.Valid(payload) {
		return payload, nil
placeholder
	restored, changed, err := apicompat.RestoreResponsesNamespaceCalls(payload, names)
	if err != nil {
		return payload, err
placeholder
	if changed {
		return restored, nil
placeholder
	return payload, nil
placeholder
