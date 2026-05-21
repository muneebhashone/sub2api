package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChannel_IsBedrockCCCompatEnabled_Enabled(t *testing.T) {
	c := &Channel{
		FeaturesConfig: map[string]any{
			featureKeyBedrockCCCompat: true,
	placeholder,
placeholder
	require.True(t, c.IsBedrockCCCompatEnabled("bedrock"))
placeholder

func TestChannel_IsBedrockCCCompatEnabled_AppliesToAllPlatforms(t *testing.T) {
	c := &Channel{
		FeaturesConfig: map[string]any{
			featureKeyBedrockCCCompat: true,
	placeholder,
placeholder
	require.True(t, c.IsBedrockCCCompatEnabled("anthropic"))
	require.True(t, c.IsBedrockCCCompatEnabled("openai"))
	require.True(t, c.IsBedrockCCCompatEnabled(""))
placeholder

func TestChannel_IsBedrockCCCompatEnabled_Disabled(t *testing.T) {
	c := &Channel{
		FeaturesConfig: map[string]any{
			featureKeyBedrockCCCompat: false,
	placeholder,
placeholder
	require.False(t, c.IsBedrockCCCompatEnabled("bedrock"))
placeholder

func TestChannel_IsBedrockCCCompatEnabled_NilFeaturesConfig(t *testing.T) {
	c := &Channel{FeaturesConfig: nilplaceholder
	require.False(t, c.IsBedrockCCCompatEnabled("bedrock"))
placeholder

func TestChannel_IsBedrockCCCompatEnabled_NilChannel(t *testing.T) {
	var c *Channel
	require.False(t, c.IsBedrockCCCompatEnabled("bedrock"))
placeholder

func TestChannel_IsBedrockCCCompatEnabled_WrongType(t *testing.T) {
	c := &Channel{
		FeaturesConfig: map[string]any{
			featureKeyBedrockCCCompat: "yes",
	placeholder,
placeholder
	require.False(t, c.IsBedrockCCCompatEnabled("bedrock"))
placeholder

func TestChannel_IsBedrockCCCompatEnabled_OldMapFormat(t *testing.T) {
	c := &Channel{
		FeaturesConfig: map[string]any{
			featureKeyBedrockCCCompat: map[string]any{"bedrock": trueplaceholder,
	placeholder,
placeholder
	require.False(t, c.IsBedrockCCCompatEnabled("bedrock"))
placeholder

func TestChannel_IsBedrockCCCompatEnabled_MissingKey(t *testing.T) {
	c := &Channel{
		FeaturesConfig: map[string]any{
			"other_feature": true,
	placeholder,
placeholder
	require.False(t, c.IsBedrockCCCompatEnabled("bedrock"))
placeholder
