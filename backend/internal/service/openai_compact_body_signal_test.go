//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
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
	require.True(t, hasCompactionTriggerInInput(body))
placeholder

func TestHasCompactionTriggerInInput_NoTrigger(t *testing.T) {
	body := []byte(`{
		"model":"gpt-5.5",
		"input":[
			{"type":"message","role":"user","content":"hello"placeholder
		]
placeholder`)
	require.False(t, hasCompactionTriggerInInput(body))
placeholder

func TestHasCompactionTriggerInInput_EmptyInput(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","input":[]placeholder`)
	require.False(t, hasCompactionTriggerInInput(body))
placeholder

func TestHasCompactionTriggerInInput_NoInputField(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5"placeholder`)
	require.False(t, hasCompactionTriggerInInput(body))
placeholder

func TestHasCompactionTriggerInInput_EmptyBody(t *testing.T) {
	require.False(t, hasCompactionTriggerInInput(nil))
	require.False(t, hasCompactionTriggerInInput([]byte{placeholder))
placeholder

func TestHasCompactionTriggerInInput_StringInput(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","input":"compaction_trigger"placeholder`)
	require.False(t, hasCompactionTriggerInInput(body))
placeholder

func TestHasCompactionTriggerInInput_CompactTriggerOnly(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","input":[{"type":"compaction_trigger"placeholder]placeholder`)
	require.True(t, hasCompactionTriggerInInput(body))
placeholder
