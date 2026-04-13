//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetWebSearchEmulationMode_Enabled(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{featureKeyWebSearchEmulation: "enabled"placeholder,
placeholder
	require.Equal(t, WebSearchModeEnabled, a.GetWebSearchEmulationMode())
placeholder

func TestGetWebSearchEmulationMode_Disabled(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{featureKeyWebSearchEmulation: "disabled"placeholder,
placeholder
	require.Equal(t, WebSearchModeDisabled, a.GetWebSearchEmulationMode())
placeholder

func TestGetWebSearchEmulationMode_Default(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{featureKeyWebSearchEmulation: "default"placeholder,
placeholder
	require.Equal(t, WebSearchModeDefault, a.GetWebSearchEmulationMode())
placeholder

func TestGetWebSearchEmulationMode_UnknownString(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{featureKeyWebSearchEmulation: "unknown"placeholder,
placeholder
	require.Equal(t, WebSearchModeDefault, a.GetWebSearchEmulationMode())
placeholder

func TestGetWebSearchEmulationMode_OldBoolTrue(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{featureKeyWebSearchEmulation: trueplaceholder,
placeholder
	// bool is not a string, type assertion fails → default
	require.Equal(t, WebSearchModeDefault, a.GetWebSearchEmulationMode())
placeholder

func TestGetWebSearchEmulationMode_OldBoolFalse(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{featureKeyWebSearchEmulation: falseplaceholder,
placeholder
	require.Equal(t, WebSearchModeDefault, a.GetWebSearchEmulationMode())
placeholder

func TestGetWebSearchEmulationMode_NilAccount(t *testing.T) {
	var a *Account
	require.Equal(t, WebSearchModeDefault, a.GetWebSearchEmulationMode())
placeholder

func TestGetWebSearchEmulationMode_NilExtra(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Extra:    nil,
placeholder
	require.Equal(t, WebSearchModeDefault, a.GetWebSearchEmulationMode())
placeholder

func TestGetWebSearchEmulationMode_MissingField(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{placeholder,
placeholder
	require.Equal(t, WebSearchModeDefault, a.GetWebSearchEmulationMode())
placeholder

func TestGetWebSearchEmulationMode_NonAnthropicPlatform(t *testing.T) {
	a := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{featureKeyWebSearchEmulation: "enabled"placeholder,
placeholder
	require.Equal(t, WebSearchModeDefault, a.GetWebSearchEmulationMode())
placeholder

func TestGetWebSearchEmulationMode_NonAPIKeyType(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeOAuth,
		Extra:    map[string]any{featureKeyWebSearchEmulation: "enabled"placeholder,
placeholder
	require.Equal(t, WebSearchModeDefault, a.GetWebSearchEmulationMode())
placeholder
