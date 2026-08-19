//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestHasCompactionTriggerInInput_DetectsCompactSignal(t *testing.T) {
	body := []byte(`{
		"model":"gpt-5.5",
		"stream":true,
		"input":[
			{"type":"message","role":"user","content":"hello"placeholder,
			{"type":"compaction_trigger"placeholder
		]
placeholder`)
	require.True(t, HasCompactionTriggerInInput(body))
placeholder

func TestHasCompactionTriggerInInput_NoTrigger(t *testing.T) {
	body := []byte(`{
		"model":"gpt-5.5",
		"input":[
			{"type":"message","role":"user","content":"hello"placeholder
		]
placeholder`)
	require.False(t, HasCompactionTriggerInInput(body))
placeholder

func TestHasCompactionTriggerInInput_EmptyInput(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","input":[]placeholder`)
	require.False(t, HasCompactionTriggerInInput(body))
placeholder

func TestHasCompactionTriggerInInput_NoInputField(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5"placeholder`)
	require.False(t, HasCompactionTriggerInInput(body))
placeholder

func TestHasCompactionTriggerInInput_EmptyBody(t *testing.T) {
	require.False(t, HasCompactionTriggerInInput(nil))
	require.False(t, HasCompactionTriggerInInput([]byte{placeholder))
placeholder

func TestHasCompactionTriggerInInput_StringInput(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","input":"compaction_trigger"placeholder`)
	require.False(t, HasCompactionTriggerInInput(body))
placeholder

func TestHasCompactionTriggerInInput_CompactTriggerOnly(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","input":[{"type":"compaction_trigger"placeholder]placeholder`)
	require.True(t, HasCompactionTriggerInInput(body))
placeholder

func TestNormalizeCompactionTriggerInputOrder_MovesAndCollapsesTriggers(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","input":[{"type":"compaction_trigger"placeholder,{"type":"message","role":"user","content":"tail"placeholder,{"type":"compaction_trigger"placeholder]placeholder`)
	normalized, changed, err := NormalizeCompactionTriggerInputOrder(body)
placeholder
	require.True(t, changed)
	items := gjson.GetBytes(normalized, "input").Array()
	require.Len(t, items, 2)
	require.Equal(t, "message", items[0].Get("type").String())
	require.Equal(t, "compaction_trigger", items[1].Get("type").String())
placeholder

func TestNormalizeCompactionTriggerInputOrder_AlreadyFinalPreservesBytes(t *testing.T) {
	body := []byte(`{"input":[{"type":"message"placeholder,{"type":"compaction_trigger"placeholder]placeholder`)
	normalized, changed, err := NormalizeCompactionTriggerInputOrder(body)
placeholder
	require.False(t, changed)
	require.Equal(t, string(body), string(normalized))
placeholder
