package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

const grokResponsesClientToolMappingContextKey = "grok_responses_client_tool_mapping"

func adaptResponsesClientToolsForFunctionUpstream(body []byte, upstream string) ([]byte, apicompat.ResponsesClientToolMapping, error) {
	return adaptResponsesClientToolsForFunctionUpstreamWithMapping(
		body,
		upstream,
		apicompat.ResponsesClientToolMapping{placeholder,
	)
placeholder

func adaptResponsesClientToolsForFunctionUpstreamWithMapping(
	body []byte,
	upstream string,
	inherited apicompat.ResponsesClientToolMapping,
) ([]byte, apicompat.ResponsesClientToolMapping, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	var requestBody map[string]any
	if err := decoder.Decode(&requestBody); err != nil {
		return body, apicompat.ResponsesClientToolMapping{placeholder, fmt.Errorf("decode %s Responses client tools: %w", upstream, err)
placeholder

	mapping, changed, err := apicompat.AdaptResponsesClientToolsWithInheritedMapping(requestBody, inherited)
	if err != nil {
		return body, apicompat.ResponsesClientToolMapping{placeholder, err
placeholder
	if !changed {
		return body, mapping, nil
placeholder
	rebuilt, err := marshalOpenAIUpstreamJSON(requestBody)
	if err != nil {
		return body, apicompat.ResponsesClientToolMapping{placeholder, fmt.Errorf("encode %s Responses client tools: %w", upstream, err)
placeholder
	return rebuilt, mapping, nil
placeholder

func adaptGrokResponsesClientTools(body []byte) ([]byte, apicompat.ResponsesClientToolMapping, error) {
	return adaptResponsesClientToolsForFunctionUpstream(body, "Grok")
placeholder

func hasResponsesClientToolMapping(mapping apicompat.ResponsesClientToolMapping) bool {
	return len(mapping.CustomTools) > 0 || mapping.ToolSearch || len(mapping.NamespaceTools) > 0
placeholder

func hasGrokResponsesClientToolMapping(mapping apicompat.ResponsesClientToolMapping) bool {
	return hasResponsesClientToolMapping(mapping)
placeholder

func setGrokResponsesClientToolMapping(c *gin.Context, mapping apicompat.ResponsesClientToolMapping) {
	if c == nil {
		return
placeholder
	if !hasGrokResponsesClientToolMapping(mapping) {
		clearGrokResponsesClientToolMapping(c)
		return
placeholder
	c.Set(grokResponsesClientToolMappingContextKey, mapping)
placeholder

func clearGrokResponsesClientToolMapping(c *gin.Context) {
	if c == nil {
		return
placeholder
	if _, exists := c.Get(grokResponsesClientToolMappingContextKey); !exists {
		return
placeholder
	c.Set(grokResponsesClientToolMappingContextKey, apicompat.ResponsesClientToolMapping{placeholder)
placeholder

func grokResponsesClientToolMapping(c *gin.Context) (apicompat.ResponsesClientToolMapping, bool) {
	if c == nil {
		return apicompat.ResponsesClientToolMapping{placeholder, false
placeholder
	value, ok := c.Get(grokResponsesClientToolMappingContextKey)
	if !ok {
		return apicompat.ResponsesClientToolMapping{placeholder, false
placeholder
	mapping, ok := value.(apicompat.ResponsesClientToolMapping)
	return mapping, ok && hasGrokResponsesClientToolMapping(mapping)
placeholder

func restoreGrokResponsesClientToolPayload(c *gin.Context, payload []byte) ([]byte, error) {
	mapping, ok := grokResponsesClientToolMapping(c)
	if !ok || !bytes.Contains(payload, []byte(`"function_call"`)) || !json.Valid(payload) {
		return payload, nil
placeholder
	restored, _, err := apicompat.RestoreResponsesClientToolPayload(payload, mapping)
	return restored, err
placeholder

type responsesClientToolStreamBody struct {
	*io.PipeReader
	source io.Closer
placeholder

func (b *responsesClientToolStreamBody) Close() error {
	readerErr := b.PipeReader.Close()
	sourceErr := b.source.Close()
	if readerErr != nil {
		return readerErr
placeholder
	return sourceErr
placeholder

func newResponsesClientToolStreamBody(
	source io.ReadCloser,
	mapping apicompat.ResponsesClientToolMapping,
	maxLineSize int,
) io.ReadCloser {
	reader, writer := io.Pipe()
	body := &responsesClientToolStreamBody{PipeReader: reader, source: sourceplaceholder
	go transformResponsesClientToolStream(source, writer, mapping, maxLineSize)
	return body
placeholder

func newGrokResponsesClientToolStreamBody(
	source io.ReadCloser,
	mapping apicompat.ResponsesClientToolMapping,
	maxLineSize int,
) io.ReadCloser {
	return newResponsesClientToolStreamBody(source, mapping, maxLineSize)
placeholder

func transformResponsesClientToolStream(
	source io.ReadCloser,
	destination *io.PipeWriter,
	mapping apicompat.ResponsesClientToolMapping,
	maxLineSize int,
) {
	defer func() { _ = source.Close() placeholder()
	if maxLineSize <= 0 {
		maxLineSize = defaultMaxLineSize
placeholder

	scanner := bufio.NewScanner(source)
	scanBuf := getSSEScannerBuf64K()
	defer putSSEScannerBuf64K(scanBuf)
	scanner.Buffer(scanBuf[:0], maxLineSize)
	documents := newOpenAISSEJSONDocumentScanner(scanner)
	restorer := apicompat.NewResponsesClientToolStreamRestorer(mapping)
	buffered := bufio.NewWriterSize(destination, 4*1024)
	pendingFields := make([]string, 0, 2)
	frameHadEventField := false
	frameEmitted := false

	writeLine := func(line string) error {
		if _, err := buffered.WriteString(line); err != nil {
			return err
	placeholder
		return buffered.WriteByte('\n')
placeholder
	writePendingFields := func(payload []byte, includeNonEvent bool) error {
		eventType := strings.TrimSpace(gjson.GetBytes(payload, "type").String())
		for _, field := range pendingFields {
			if _, isEvent := extractOpenAISSEEventLine(field); isEvent {
				if eventType != "" {
					if err := writeLine("event: " + eventType); err != nil {
						return err
				placeholder
			placeholder else if err := writeLine(field); err != nil {
					return err
			placeholder
				continue
		placeholder
			if includeNonEvent {
				if err := writeLine(field); err != nil {
					return err
			placeholder
		placeholder
	placeholder
		return nil
placeholder
	writePayloads := func(payloads [][]byte) error {
		for index, payload := range payloads {
			if index == 0 {
				if err := writePendingFields(payload, true); err != nil {
					return err
			placeholder
		placeholder else if frameHadEventField {
				eventType := strings.TrimSpace(gjson.GetBytes(payload, "type").String())
				if eventType != "" {
					if err := writeLine("event: " + eventType); err != nil {
						return err
				placeholder
			placeholder
		placeholder
			if err := writeLine("data: " + string(payload)); err != nil {
				return err
		placeholder
			if err := writeLine(""); err != nil {
				return err
		placeholder
	placeholder
		return buffered.Flush()
placeholder

	for documents.Scan() {
		line := documents.Text()
		data, isData := extractOpenAISSEDataLine(line)
		if isData {
			payload := []byte(data)
			payloads := [][]byte{payloadplaceholder
			if json.Valid(payload) {
				var err error
				payloads, _, err = restorer.RestoreEvent(payload)
				if err != nil {
					_ = buffered.Flush()
					_ = destination.CloseWithError(fmt.Errorf("restore Responses client tool event: %w", err))
					return
			placeholder
		placeholder
			if err := writePayloads(payloads); err != nil {
				_ = destination.CloseWithError(err)
				return
		placeholder
			pendingFields = pendingFields[:0]
			frameHadEventField = false
			frameEmitted = true
			continue
	placeholder

		if line == "" {
			if !frameEmitted {
				for _, field := range pendingFields {
					if err := writeLine(field); err != nil {
						_ = destination.CloseWithError(err)
						return
				placeholder
			placeholder
				if len(pendingFields) > 0 {
					if err := writeLine(""); err != nil {
						_ = destination.CloseWithError(err)
						return
				placeholder
					if err := buffered.Flush(); err != nil {
						_ = destination.CloseWithError(err)
						return
				placeholder
			placeholder
		placeholder
			pendingFields = pendingFields[:0]
			frameHadEventField = false
			frameEmitted = false
			continue
	placeholder

		if _, isEvent := extractOpenAISSEEventLine(line); isEvent {
			frameHadEventField = true
	placeholder
		pendingFields = append(pendingFields, line)
placeholder

	for _, field := range pendingFields {
		if err := writeLine(field); err != nil {
			_ = destination.CloseWithError(err)
			return
	placeholder
placeholder
	if err := buffered.Flush(); err != nil {
		_ = destination.CloseWithError(err)
		return
placeholder
	if err := documents.Err(); err != nil {
		_ = destination.CloseWithError(err)
		return
placeholder
	_ = destination.Close()
placeholder
