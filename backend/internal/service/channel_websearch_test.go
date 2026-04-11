package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChannel_IsWebSearchEmulationEnabled_Enabled(t *testing.T) {
	c := &Channel{
		FeaturesConfig: map[string]any{
			featureKeyWebSearchEmulation: map[string]any{"anthropic": trueplaceholder,
	placeholder,
placeholder
	require.True(t, c.IsWebSearchEmulationEnabled("anthropic"))
placeholder

func TestChannel_IsWebSearchEmulationEnabled_DifferentPlatform(t *testing.T) {
	c := &Channel{
		FeaturesConfig: map[string]any{
			featureKeyWebSearchEmulation: map[string]any{"anthropic": trueplaceholder,
	placeholder,
placeholder
	require.False(t, c.IsWebSearchEmulationEnabled("openai"))
placeholder

func TestChannel_IsWebSearchEmulationEnabled_Disabled(t *testing.T) {
	c := &Channel{
		FeaturesConfig: map[string]any{
			featureKeyWebSearchEmulation: map[string]any{"anthropic": falseplaceholder,
	placeholder,
placeholder
	require.False(t, c.IsWebSearchEmulationEnabled("anthropic"))
placeholder

func TestChannel_IsWebSearchEmulationEnabled_NilFeaturesConfig(t *testing.T) {
	c := &Channel{FeaturesConfig: nilplaceholder
	require.False(t, c.IsWebSearchEmulationEnabled("anthropic"))
placeholder

func TestChannel_IsWebSearchEmulationEnabled_NilChannel(t *testing.T) {
	var c *Channel
	require.False(t, c.IsWebSearchEmulationEnabled("anthropic"))
placeholder

func TestChannel_IsWebSearchEmulationEnabled_WrongStructure(t *testing.T) {
	c := &Channel{
		FeaturesConfig: map[string]any{
			featureKeyWebSearchEmulation: true, // not a map
	placeholder,
placeholder
	require.False(t, c.IsWebSearchEmulationEnabled("anthropic"))
placeholder

func TestChannel_IsWebSearchEmulationEnabled_PlatformValueNotBool(t *testing.T) {
	c := &Channel{
		FeaturesConfig: map[string]any{
			featureKeyWebSearchEmulation: map[string]any{"anthropic": "yes"placeholder,
	placeholder,
placeholder
	require.False(t, c.IsWebSearchEmulationEnabled("anthropic"))
placeholder
