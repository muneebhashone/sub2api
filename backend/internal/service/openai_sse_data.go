package service

import (
	"strings"

	"github.com/tidwall/gjson"
)

type openAISSEDataAccumulator struct {
	lines []string
placeholder

func (a *openAISSEDataAccumulator) AddLine(line string, fn func([]byte)) {
	if fn == nil {
		return
placeholder
	trimmedLine := strings.TrimRight(line, "\r\n")
	if data, ok := extractOpenAISSEDataLine(trimmedLine); ok {
		a.lines = append(a.lines, data)
		return
placeholder
	if strings.TrimSpace(trimmedLine) == "" {
		a.Flush(fn)
placeholder
placeholder

func (a *openAISSEDataAccumulator) Flush(fn func([]byte)) {
	if fn == nil || len(a.lines) == 0 {
		return
placeholder
	emitOpenAISSEDataPayloads(a.lines, fn)
	a.lines = a.lines[:0]
placeholder

func forEachOpenAISSEDataPayload(body string, fn func([]byte)) {
	if fn == nil || strings.TrimSpace(body) == "" {
		return
placeholder
	var acc openAISSEDataAccumulator
	for _, line := range strings.Split(body, "\n") {
		acc.AddLine(line, fn)
placeholder
	acc.Flush(fn)
placeholder

func emitOpenAISSEDataPayloads(lines []string, fn func([]byte)) {
	if fn == nil || len(lines) == 0 {
		return
placeholder
	if len(lines) == 1 {
		emitOpenAISSEDataPayload(lines[0], fn)
		return
placeholder
	joined := strings.Join(lines, "\n")
	if gjson.Valid(joined) {
		emitOpenAISSEDataPayload(joined, fn)
		return
placeholder
	for _, line := range lines {
		emitOpenAISSEDataPayload(line, fn)
placeholder
placeholder

func emitOpenAISSEDataPayload(data string, fn func([]byte)) {
	data = strings.TrimSpace(data)
	if data == "" || data == "[DONE]" {
		return
placeholder
	fn([]byte(data))
placeholder
