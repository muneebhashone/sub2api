package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
)

type grokResponsesBillingPingFilterBody struct {
	*io.PipeReader
	source    io.Closer
	closeOnce sync.Once
	closeErr  error
placeholder

func (b *grokResponsesBillingPingFilterBody) Close() error {
	readerErr := b.PipeReader.Close()
	sourceErr := b.closeSource()
	if readerErr != nil {
		return readerErr
placeholder
	return sourceErr
placeholder

func (b *grokResponsesBillingPingFilterBody) closeSource() error {
	b.closeOnce.Do(func() { b.closeErr = b.source.Close() placeholder)
	return b.closeErr
placeholder

func newGrokResponsesBillingPingFilterBody(source io.ReadCloser, account *Account, maxLineSize int) io.ReadCloser {
	if account == nil || account.Platform != PlatformGrok {
		return source
placeholder
	reader, writer := io.Pipe()
	body := &grokResponsesBillingPingFilterBody{PipeReader: reader, source: sourceplaceholder
	go filterGrokResponsesBillingPings(source, writer, body.closeSource, maxLineSize)
	return body
placeholder

func filterGrokResponsesBillingPings(
	source io.Reader,
	destination *io.PipeWriter,
	closeSource func() error,
	maxLineSize int,
) {
	defer func() { _ = closeSource() placeholder()
	if maxLineSize <= 0 {
		maxLineSize = defaultMaxLineSize
placeholder

	scanner := bufio.NewScanner(source)
	scanBuf := getSSEScannerBuf64K()
	defer putSSEScannerBuf64K(scanBuf)
	initialBufferSize := len(scanBuf)
	if maxLineSize < initialBufferSize {
		initialBufferSize = maxLineSize
placeholder
	scanner.Buffer(scanBuf[:0:initialBufferSize], maxLineSize)
	scanner.Split(scanSSELinesPreservingEndings)
	frame := make([][]byte, 0, 3)

	writeFrame := func() error {
		if !isGrokResponsesBillingPingFrame(frame) {
			for _, line := range frame {
				if _, err := destination.Write(line); err != nil {
					return err
			placeholder
		placeholder
	placeholder
		frame = frame[:0]
		return nil
placeholder

	for scanner.Scan() {
		line := append([]byte(nil), scanner.Bytes()...)
		frame = append(frame, line)
		if len(bytes.TrimSuffix(bytes.TrimSuffix(line, []byte("\n")), []byte("\r"))) == 0 {
			if err := writeFrame(); err != nil {
				_ = destination.CloseWithError(err)
				return
		placeholder
	placeholder
placeholder
	if len(frame) > 0 {
		if err := writeFrame(); err != nil {
			_ = destination.CloseWithError(err)
			return
	placeholder
placeholder
	if err := scanner.Err(); err != nil {
		_ = destination.CloseWithError(fmt.Errorf("filter Grok Responses billing ping: %w", err))
		return
placeholder
	_ = destination.Close()
placeholder

func scanSSELinesPreservingEndings(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for index, value := range data {
		switch value {
		case '\n':
			return index + 1, data[:index+1], nil
		case '\r':
			if index+1 == len(data) && !atEOF {
				return 0, nil, nil
		placeholder
			if index+1 < len(data) && data[index+1] == '\n' {
				return index + 2, data[:index+2], nil
		placeholder
			return index + 1, data[:index+1], nil
	placeholder
placeholder
	if atEOF && len(data) > 0 {
		return len(data), data, nil
placeholder
	return 0, nil, nil
placeholder

func isGrokResponsesBillingPingFrame(rawLines [][]byte) bool {
	eventType := ""
	data := ""
	eventLines := 0
	dataLines := 0
	for _, rawLine := range rawLines {
		line := strings.TrimSuffix(strings.TrimSuffix(string(rawLine), "\n"), "\r")
		if line == "" {
			continue
	placeholder
		if value, ok := extractOpenAISSEEventLine(line); ok {
			eventType = value
			eventLines++
			continue
	placeholder
		if value, ok := extractOpenAISSEDataLine(line); ok {
			data = value
			dataLines++
			continue
	placeholder
		return false
placeholder
	if eventType != "ping" || eventLines != 1 || dataLines != 1 {
		return false
placeholder

	var payload struct {
		Type         string          `json:"type"`
		OpenCodeType json.RawMessage `json:"x-opencode-type"`
		Cost         json.RawMessage `json:"cost"`
placeholder
	if err := json.Unmarshal([]byte(data), &payload); err != nil || payload.Type != "ping" {
		return false
placeholder
	if len(payload.OpenCodeType) > 0 {
		var openCodeType string
		return json.Unmarshal(payload.OpenCodeType, &openCodeType) == nil && openCodeType == "inference-cost"
placeholder
	var cost string
	return json.Unmarshal(payload.Cost, &cost) == nil && cost == "0"
placeholder
