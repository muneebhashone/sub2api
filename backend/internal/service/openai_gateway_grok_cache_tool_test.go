package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestAppendMissingGrokFreeCacheNativeTools_PureClientFunctionNoInject(t *testing.T) {
	body := []byte(`{
		"model": "grok-4.5",
		"tools": [
			{"type":"function","name":"view_image","description":"View image","parameters":{"type":"object","properties":{"path":{"type":"string"placeholderplaceholder,"required":["path"]placeholderplaceholder
		],
		"tool_choice": "auto"
placeholder`)

	result, err := appendMissingGrokFreeCacheNativeTools(body)
placeholder

	tools := gjson.GetBytes(result, "tools").Array()
	for _, tool := range tools {
		toolType := tool.Get("type").String()
		assert.NotEqual(t, "web_search", toolType, "should not inject web_search for pure client functions")
		assert.NotEqual(t, "x_search", toolType, "should not inject x_search for pure client functions")
placeholder
placeholder

func TestAppendMissingGrokFreeCacheNativeTools_FunctionPlusWebSearchInjects(t *testing.T) {
	body := []byte(`{
		"model": "grok-4.5",
		"tools": [
			{"type":"function","name":"view_image","description":"View","parameters":{"type":"object","properties":{"path":{"type":"string"placeholderplaceholder,"required":["path"]placeholderplaceholder,
			{"type":"function","name":"web_search","description":"Search","parameters":{"type":"object","properties":{"query":{"type":"string"placeholderplaceholder,"required":["query"]placeholderplaceholder
		]
placeholder`)

	result, err := appendMissingGrokFreeCacheNativeTools(body)
placeholder

	tools := gjson.GetBytes(result, "tools").Array()
	types := make(map[string]bool)
	for _, tool := range tools {
		types[tool.Get("type").String()] = true
placeholder
	assert.True(t, types["web_search"], "web_search should be present (converted from function)")
	assert.True(t, types["x_search"], "x_search should be injected when web_search is present alongside client functions")
placeholder

func TestAppendMissingGrokFreeCacheNativeTools_NativeSearchAlreadyPresent(t *testing.T) {
	body := []byte(`{
		"model": "grok-4.5",
		"tools": [
			{"type":"function","name":"view_image","description":"View","parameters":{"type":"object","properties":{"path":{"type":"string"placeholderplaceholder,"required":["path"]placeholderplaceholder,
			{"type":"web_search"placeholder
		]
placeholder`)

	result, err := appendMissingGrokFreeCacheNativeTools(body)
placeholder

	tools := gjson.GetBytes(result, "tools").Array()
	types := make(map[string]bool)
	for _, tool := range tools {
		types[tool.Get("type").String()] = true
placeholder
	assert.True(t, types["web_search"])
	assert.True(t, types["x_search"], "x_search should be injected when web_search is already present")
placeholder

func TestAppendMissingGrokFreeCacheNativeTools_MultipleFunctionsNoSearch(t *testing.T) {
	body := []byte(`{
		"model": "grok-4.5",
		"tools": [
			{"type":"function","name":"view_image","description":"View","parameters":{"type":"object","properties":{"path":{"type":"string"placeholderplaceholder,"required":["path"]placeholderplaceholder,
			{"type":"function","name":"read_file","description":"Read","parameters":{"type":"object","properties":{"path":{"type":"string"placeholderplaceholder,"required":["path"]placeholderplaceholder
		]
placeholder`)

	result, err := appendMissingGrokFreeCacheNativeTools(body)
placeholder

	tools := gjson.GetBytes(result, "tools").Array()
	require.Len(t, tools, 2, "no tools should be injected for pure client functions")
placeholder
