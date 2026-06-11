package service

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/tidwall/gjson"
)

type openAIImageOutputCounter struct {
	seen         map[string]struct{placeholder
	seenSizes    map[string]string
	seenOrder    []string
	dataSizes    []string
	count        int
	maxDataCount int
placeholder

func newOpenAIImageOutputCounter() *openAIImageOutputCounter {
	return &openAIImageOutputCounter{
		seen:      make(map[string]struct{placeholder),
		seenSizes: make(map[string]string),
placeholder
placeholder

func (c *openAIImageOutputCounter) Count() int {
	if c == nil {
		return 0
placeholder
	if c.maxDataCount > c.count {
		return c.maxDataCount
placeholder
	return c.count
placeholder

func (c *openAIImageOutputCounter) Sizes() []string {
	if c == nil {
		return nil
placeholder
	sizes := make([]string, 0, len(c.seenOrder)+len(c.dataSizes))
	for _, key := range c.seenOrder {
		if size := strings.TrimSpace(c.seenSizes[key]); size != "" {
			sizes = append(sizes, size)
	placeholder
placeholder
	if len(sizes) == 0 && len(c.dataSizes) > 0 {
		sizes = append(sizes, c.dataSizes...)
placeholder
	if len(sizes) == 0 {
		return nil
placeholder
	return sizes
placeholder

func (c *openAIImageOutputCounter) AddJSONResponse(body []byte) {
	if c == nil || len(body) == 0 || !gjson.ValidBytes(body) {
		return
placeholder
	c.addDataArray(gjson.GetBytes(body, "data"))
	c.addOutputArray(gjson.GetBytes(body, "output"))
	c.addOutputArray(gjson.GetBytes(body, "response.output"))
placeholder

func (c *openAIImageOutputCounter) AddSSEData(data []byte) {
	if c == nil || len(data) == 0 || strings.TrimSpace(string(data)) == "[DONE]" || !gjson.ValidBytes(data) {
		return
placeholder
	root := gjson.ParseBytes(data)
	c.addDataArray(root.Get("data"))
	eventType := strings.TrimSpace(root.Get("type").String())
	switch eventType {
	case "response.output_item.done":
		c.addImageOutputItem(root.Get("item"))
	case "response.completed", "response.done":
		c.addOutputArray(root.Get("response.output"))
	case "image_generation.completed":
		if item := root.Get("item"); item.Exists() {
			c.addImageOutputItem(item)
			return
	placeholder
		if output := root.Get("output"); output.Exists() {
			c.addImageOutputItem(output)
			return
	placeholder
		c.addImageOutputItem(root)
placeholder
placeholder

func (c *openAIImageOutputCounter) AddSSEBody(body string) {
	if c == nil || strings.TrimSpace(body) == "" {
		return
placeholder
	forEachOpenAISSEDataPayload(body, c.AddSSEData)
placeholder

func (c *openAIImageOutputCounter) addDataArray(data gjson.Result) {
	if !data.IsArray() {
		return
placeholder
	items := data.Array()
	imageCount := 0
	sizes := make([]string, 0, len(items))
	for _, item := range items {
		if !item.IsObject() {
			continue
	placeholder
		hasImageOutput := strings.TrimSpace(item.Get("url").String()) != "" ||
			strings.TrimSpace(item.Get("b64_json").String()) != ""
		if !hasImageOutput {
			continue
	placeholder
		imageCount++
		if size := strings.TrimSpace(item.Get("size").String()); size != "" {
			sizes = append(sizes, size)
	placeholder
placeholder
	if imageCount > c.maxDataCount {
		c.maxDataCount = imageCount
placeholder
	if len(sizes) > 0 {
		c.dataSizes = sizes
placeholder
placeholder

func (c *openAIImageOutputCounter) addOutputArray(output gjson.Result) {
	if !output.IsArray() {
		return
placeholder
	output.ForEach(func(_, item gjson.Result) bool {
		c.addImageOutputItem(item)
		return true
placeholder)
placeholder

func (c *openAIImageOutputCounter) addImageOutputItem(item gjson.Result) {
	if !item.Exists() || !item.IsObject() {
		return
placeholder
	itemType := strings.TrimSpace(item.Get("type").String())
	if itemType != "" && itemType != "image_generation_call" && itemType != "image_generation.completed" {
		return
placeholder
	if strings.Contains(strings.ToLower(item.Raw), "partial_image") {
		return
placeholder
	result := strings.TrimSpace(item.Get("result").String())
	if result == "" {
		result = strings.TrimSpace(item.Get("b64_json").String())
placeholder
	if result == "" {
		result = strings.TrimSpace(item.Get("url").String())
placeholder
	if result == "" {
		return
placeholder
	key := strings.TrimSpace(item.Get("id").String())
	if key == "" {
		key = strings.TrimSpace(item.Get("call_id").String())
placeholder
	if key == "" {
		key = hashOpenAIImageOutputResult(result)
placeholder
	if key == "" {
		return
placeholder
	size := strings.TrimSpace(item.Get("size").String())
	if _, exists := c.seen[key]; exists {
		if size != "" && strings.TrimSpace(c.seenSizes[key]) == "" {
			c.seenSizes[key] = size
	placeholder
		return
placeholder
	c.seen[key] = struct{placeholder{placeholder
	c.seenOrder = append(c.seenOrder, key)
	if size != "" {
		c.seenSizes[key] = size
placeholder
	c.count++
placeholder

func hashOpenAIImageOutputResult(result string) string {
	result = strings.TrimSpace(result)
	if result == "" {
		return ""
placeholder
	sum := sha256.Sum256([]byte(result))
	return hex.EncodeToString(sum[:])
placeholder

func countOpenAIResponseImageOutputsFromJSONBytes(body []byte) int {
	counter := newOpenAIImageOutputCounter()
	counter.AddJSONResponse(body)
	return counter.Count()
placeholder

func collectOpenAIResponseImageOutputSizesFromJSONBytes(body []byte) []string {
	counter := newOpenAIImageOutputCounter()
	counter.AddJSONResponse(body)
	return counter.Sizes()
placeholder

func countOpenAIImageOutputsFromSSEBody(body string) int {
	counter := newOpenAIImageOutputCounter()
	counter.AddSSEBody(body)
	return counter.Count()
placeholder

func collectOpenAIImageOutputSizesFromSSEBody(body string) []string {
	counter := newOpenAIImageOutputCounter()
	counter.AddSSEBody(body)
	return counter.Sizes()
placeholder
